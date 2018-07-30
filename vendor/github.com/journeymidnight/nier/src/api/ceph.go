package api

import (
	"github.com/journeymidnight/go-ceph/rados"
	"github.com/journeymidnight/nier/src/helper"
	"encoding/json"
	"bytes"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	. "github.com/journeymidnight/nier/src/api/datatype"
	. "github.com/journeymidnight/nier/src/error"
	"os/exec"
)
//import "fmt"

//for snmp VariableTypeCounter64
//copied from https://github.com/digitalocean/ceph_exporter

type DiskStatus struct {
	Name               string
	UsedKb             int64
	Kb                 int64
	Up                 int64
	In                 int64
	Host               string
	HostClusterAddress string //ceph osd's cluster_addr, we have only one interface
}

type CephCollector struct {
	clusterfsid      string
	ClusterStatus    string
	ClusterDisks     map[string]*DiskStatus
	ClusterUsed      int64
	ClusterAvali     int64
	clusterUpOsd     int64
	clusterInOsd     int64
	conn             *rados.Conn
	clusterMonLeader string
}

func NewCephCollector() (*CephCollector, error) {
	conn, _ := rados.NewConn()
	conn.ReadDefaultConfigFile()
	conn.Connect()
	cephCollector := &CephCollector{
		conn:          conn,
		ClusterStatus: "",
		ClusterDisks:  make(map[string]*DiskStatus),
		ClusterUsed:   0,
		ClusterAvali:  0}
	return cephCollector, nil
}

//CephCollector use C version rados, it must be shutdown
func (c *CephCollector) Shutdown() {
	c.conn.Shutdown()
}

type cephHealthStats struct {
	Fsid   string `json:"fsid"`
	Health struct {
		Summary []struct {
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
		} `json:"summary"`
		OverallStatus string `json:"overall_status"`
		Checks        map[string]struct {
			Severity string `json:"severity"`
			Summary  struct {
				Message string `json:"message"`
			} `json:"summary"`
		} `json:"checks"`
	} `json:"health"`
	OSDMap struct {
		OSDMap struct {
			NumOSDs        json.Number `json:"num_osds"`
			NumUpOSDs      json.Number `json:"num_up_osds"`
			NumInOSDs      json.Number `json:"num_in_osds"`
			NumRemappedPGs json.Number `json:"num_remapped_pgs"`
		} `json:"osdmap"`
	} `json:"osdmap"`
	PGMap struct {
		PGsByState []struct {
			StateName string      `json:"state_name"`
			Count     json.Number `json:"count"`
		} `json:"pgs_by_state"`
		NumPGs json.Number `json:"num_pgs"`
	} `json:"pgmap"`
}

type cephOSDDF struct {
	OSDNodes []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		//do not use json.Number
		Children    []int64     `json:"children"`
		CrushWeight json.Number `json:"crush_weight"`
		Depth       json.Number `json:"depth"`
		Reweight    json.Number `json:"reweight"`
		KB          json.Number `json:"kb"`
		UsedKB      json.Number `json:"kb_used"`
		AvailKB     json.Number `json:"kb_avail"`
		Utilization json.Number `json:"utilization"`
		Variance    json.Number `json:"var"`
		Pgs         json.Number `json:"pgs"`
	} `json:"nodes"`

	Summary struct {
		TotalKB      json.Number `json:"total_kb"`
		TotalUsedKB  json.Number `json:"total_kb_used"`
		TotalAvailKB json.Number `json:"total_kb_avail"`
		AverageUtil  json.Number `json:"average_utilization"`
	} `json:"summary"`
}

func (c *CephCollector) cephHealthCommand(f string) []byte {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "status",
		"format": f,
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		// should fail.
		panic(err)
	}
	return cmd
}

//run "ceph mon_status"
func (c *CephCollector) getMon() error {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "mon_status",
		"format": "json",
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		// should fail.
		panic(err)
	}

	type cephMonStatus struct {
		Name string `json:"name"`
	}
	stats := &cephMonStatus{}
	buf, _, err := c.conn.MonCommand(cmd)
	if err = json.Unmarshal(buf, stats); err != nil {
		return err
	}
	c.clusterMonLeader = stats.Name
	return nil
}

func (c *CephCollector) GetHealth() error {

	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "status",
		"format": "json",
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		// should fail.
		panic(err)
	}

	stats := &cephHealthStats{}
	buf, _, err := c.conn.MonCommand(cmd)
	if err = json.Unmarshal(buf, stats); err != nil {
		return err
	}
	c.ClusterStatus = stats.Health.OverallStatus
	c.clusterfsid = stats.Fsid

	c.clusterInOsd, err = stats.OSDMap.OSDMap.NumInOSDs.Int64()
	if err != nil {
		return err
	}

	c.clusterUpOsd, err = stats.OSDMap.OSDMap.NumUpOSDs.Int64()
	if err != nil {
		return err
	}

	return nil
}

type cephOSDDump struct {
	OSDs []struct {
		OSD         json.Number `json:"osd"`
		Up          json.Number `json:"up"`
		In          json.Number `json:"in"`
		ClusterAddr string      `json:"cluster_addr"`
	} `json:"osds"`
}

// run 'ceph osd dump'

func (c *CephCollector) GetOsdStatus() error {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd dump",
		"format": "json",
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		panic(err)
	}

	buf, _, err := c.conn.MonCommand(cmd)
	if err != nil {
		return err
	}

	osdDump := &cephOSDDump{}
	if err := json.Unmarshal(buf, osdDump); err != nil {
		return err
	}
	for _, osd := range osdDump.OSDs {
		id, err := osd.OSD.Int64()
		if err != nil {
			return err
		}
		name := fmt.Sprintf("osd.%d", id)
		osdIn, err := osd.In.Int64()
		if err != nil {
			return err
		}
		osdUp, err := osd.Up.Int64()
		if err != nil {
			return err
		}

		ip_port_list := strings.Split(osd.ClusterAddr, ":")

		var clusterIP string
		if len(ip_port_list) > 1 {
			clusterIP = ip_port_list[0]
		} else {
			clusterIP = ""
		}

		if _, ok := c.ClusterDisks[name]; ok {
			c.ClusterDisks[name].In = osdIn
			c.ClusterDisks[name].Up = osdUp
			c.ClusterDisks[name].HostClusterAddress = clusterIP
		} else {
			c.ClusterDisks[name] = &DiskStatus{In: osdIn, Up: osdUp, HostClusterAddress: clusterIP}
		}
	}
	return nil

}

//run 'ceph osd df"
func (c *CephCollector) GetDiskStatus() error {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix":        "osd df",
		"format":        "json",
		"output_method": "tree",
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		// should fail.
		panic(err)
	}

	cephOsdDf := &cephOSDDF{}
	buf, _, err := c.conn.MonCommand(cmd)
	// Workaround for Ceph Jewel after 10.2.5 produces invalid json when osd is out
	buf = bytes.Replace(buf, []byte("-nan"), []byte("0"), -1)
	if err = json.Unmarshal(buf, cephOsdDf); err != nil {
		return err
	}

	diskMap := make(map[string][]int64)

	for _, node := range cephOsdDf.OSDNodes {
		if node.Type == "host" {
			diskMap[node.Name] = node.Children
		} else if node.Type == "osd" {
			/* process osd type*/
			name := node.Name
			usedKb, err := node.UsedKB.Int64()
			if err != nil {
				return err
			}

			osdKb, err := node.KB.Int64()
			if err != nil {
				return err
			}

			//set ClusterDisks
			if _, ok := c.ClusterDisks[name]; ok {
				c.ClusterDisks[name].Name = name
				c.ClusterDisks[name].UsedKb = usedKb
				c.ClusterDisks[name].Kb = osdKb
			} else {
				c.ClusterDisks[name] = &DiskStatus{Name: name, UsedKb: usedKb, Kb: osdKb}
			}
		} else {
			/*node type is default, host, rack */
		}
	}

	/*[osd.id]=[children1, children2, child2]*/
	for hostName, childrenList := range diskMap {
		for _, childrenID := range childrenList {
			osdName := fmt.Sprintf("osd.%d", childrenID)
			c.ClusterDisks[osdName].Host = hostName
		}
	}

	totalUsedKB, err := cephOsdDf.Summary.TotalUsedKB.Int64()
	if err != nil {
		return err
	}

	c.ClusterUsed = totalUsedKB

	totalAvailKB, err := cephOsdDf.Summary.TotalAvailKB.Int64()
	if err != nil {
		return err
	}

	c.ClusterAvali = totalAvailKB
	return nil
}

func DumpCephStatus(w http.ResponseWriter, r *http.Request) {
	//parts := strings.Fields(dumpcephosdtree)
	//out, err := exec.Command(parts[0], parts[1:]...).Output()
	//if err != nil {
	//      helper.Logger.Errorln("ceph osd df tree error:", err.Error())
	//      WriteErrorResponse(w, r, err)
	//      return
	//}
	//helper.Logger.Infoln("dump ceph result:", out)
	var collector *CephCollector
	collector, _ = NewCephCollector()
	defer collector.conn.Shutdown()
	err := collector.GetOsdStatus()
	if err != nil {
		helper.Logger.Error("error %v, for getOsdStatus", err)
		return
	}
	err = collector.GetDiskStatus()
	if err != nil {
		helper.Logger.Error("error %v, for getDiskStatus", err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(collector.ClusterDisks))
}

func AddDiskToCeph(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	args := fmt.Sprintf("AddNewDisk:hostname=%s,diskname=%s", req.HostName, req.DiskName)
	cmd := exec.Command("fab", args)
	cmd.Dir = "/binaries/storedeployer"
	err = cmd.Run()
	if err != nil {
		helper.Logger.Error("faild add disk:", err, req.HostName, req.DiskName)
		WriteErrorResponse(w, r, ErrFailedAddNewDisk)
		return
	}
	WriteSuccessResponse(w, nil)
}
