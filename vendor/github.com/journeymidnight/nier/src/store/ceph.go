package store

import (
	. "github.com/journeymidnight/nier/src/api/datatype"
	. "github.com/journeymidnight/nier/src/error"
	"github.com/journeymidnight/nier/src/helper"
	"time"
	"github.com/journeymidnight/go-ceph/rados"
	"encoding/json"
	"math"
	"fmt"
	"strings"
	"sort"
	"strconv"
	"errors"
)

const (
	UserTable = "DEADBEAF_USER_TABLE"
	TokenTable = "DEADBEAF_TOKEN_TABLE"
	SnmpTable = "DEADBEAF_SNMP_TABLE"
	NtpTable = "DEADBEAF_NTP_TABLE"
	HistoryTable = "DEADBEAF_HISTORY_TABLE"
	SendedTable = "DEADBEAF_SENDED_TABLE"
	ResolvedTable = "DEADBEAF_RESOLVED_TABLE"
)

type CephStore struct{
	conn *rados.Conn
	ioctx *rados.IOContext
}

func NewCephStore() *CephStore {
	conn, err := rados.NewConn()
	conn.ReadDefaultConfigFile()
	conn.Connect()
	a := CephStore{}
	a.conn = conn
	pools, err := conn.ListPools()
	if err != nil {
		helper.Logger.Println("list pools error:", err)
		panic(err)
	}
	helper.Logger.Warnln("list pools result:", pools)

	for _, v := range pools {
		if v == "nier" {
			a.ioctx, err = conn.OpenIOContext("nier")
			if err != nil {
				helper.Logger.Println("open pool nier error:", err)
				panic(err)
			}
			a.TouchObjectIfNotExist(UserTable)
			a.TouchObjectIfNotExist(TokenTable)
			a.TouchObjectIfNotExist(SnmpTable)
			a.TouchObjectIfNotExist(NtpTable)
			a.TouchObjectIfNotExist(HistoryTable)
			a.TouchObjectIfNotExist(SendedTable)
			a.TouchObjectIfNotExist(ResolvedTable)
			a.InsertUserRecord("admin", "admin", ROLE_ROOT)
			return &a
		}
		helper.Logger.Println("list pool:", v)
	}
	err = conn.MakePool("nier")
	if err != nil {
		helper.Logger.Println("make pool nier:", err)
		panic(err)
	}
	a.ioctx, err = conn.OpenIOContext("nier")
	if err != nil {
		helper.Logger.Println("open pool nier error:", err)
		panic(err)
	}
	a.TouchObjectIfNotExist(UserTable)
	a.TouchObjectIfNotExist(TokenTable)
	a.TouchObjectIfNotExist(SnmpTable)
	a.TouchObjectIfNotExist(NtpTable)
	a.TouchObjectIfNotExist(HistoryTable)
	a.TouchObjectIfNotExist(SendedTable)
	a.TouchObjectIfNotExist(ResolvedTable)
	a.InsertUserRecord("admin", "admin", ROLE_ROOT)

	return &a
}

func (s *CephStore)TouchObjectIfNotExist(name string) {
	_, err := s.ioctx.Stat(name)
	if err != nil && err.Error() == "rados: No such file or directory" {
		err = s.ioctx.Write(name, []byte("let`s rock"), 0)
		if err != nil {
			panic(fmt.Sprintf("touch object %s failed", name))
		}
		err = s.ioctx.Truncate(name, 0)
		if err != nil {
			panic(fmt.Sprintf("truncate object %s failed", name))
		}
	}
}

func (s *CephStore)InsertUserRecord(userName string, password string, accountType string) error {
	records, err := s.ioctx.GetOmapValuesByKeys(UserTable,[]string{userName})
	if err != nil {
		return err
	}
	if len(records) > 0 {
		return ErrDuplicateAddUser
	}
	record := UserRecord{userName, password, accountType}
	data, _ := json.Marshal(record)
	err = s.ioctx.SetOmap(UserTable, map[string][]byte{
		userName:data,
	})
	if err != nil {
		helper.Logger.Errorln("Error add user", userName, password, accountType, err.Error())
		return err
	}
	return err
}

func (s *CephStore)RemoveUserRecord(userName string) error {
	err := s.ioctx.RmOmapKeys(UserTable,[]string{userName})
	if err != nil {
		helper.Logger.Errorln("Error remove user", userName, err.Error())
	}
	return err
}

func (s *CephStore)ModifyUserRecord(userName string, password string, accountType string) error {
	record := UserRecord{userName, password, accountType}
	data, _ := json.Marshal(record)
	err := s.ioctx.SetOmap(UserTable, map[string][]byte{
		userName:data,
	})
	if err != nil {
		helper.Logger.Errorln("Error remove user", userName, err.Error())
	}
	return err
}

func (s *CephStore)DescribeUserRecord(userName string) (UserRecord, error) {
	var record UserRecord
	results, err := s.ioctx.GetOmapValuesByKeys(UserTable,[]string{userName})
	if err != nil {
		return record, err
	} else if len(results) > 0 {
		err = json.Unmarshal(results[userName], &record)
		if err != nil {
			return record, err
		}
	}
	return record, nil
}

func (s *CephStore)ListUserRecords() ([]UserRecord, error) {
	var records []UserRecord = make([]UserRecord, 0)
	result, err := s.ioctx.GetOmapValues(UserTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Tokens", err)
		return records, err
	}
	for _, v := range result {
		var record UserRecord
		err = json.Unmarshal(v, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *CephStore)ValidUserAndPassword(userName string, password string) (UserRecord, error) {
	var record UserRecord
	results, err := s.ioctx.GetOmapValuesByKeys(UserTable,[]string{userName})
	if err != nil {
		return record, err
	} else if len(results) > 0 {
		err = json.Unmarshal(results[userName], &record)
		if err != nil {
			return record, err
		}
		if password == record.Password {
			return record, nil
		} else {
			return record, ErrUserOrPasswordInvalid
		}
	} else {
		return record, ErrUserOrPasswordInvalid
	}

	return record, nil
}

func (s *CephStore)CheckUserExist(userName string) (bool, error) {
	results, err := s.ioctx.GetOmapValuesByKeys(UserTable,[]string{userName})
	if err != nil {
		return false, err
	}
	if len(results) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *CephStore)InsertTokenRecord(token string, userName string, userType string) error {
	created := time.Now().Format(TimeFormat)
	expired := time.Now().Add(time.Duration(helper.Config.TokenExpire * 1000000000)).Format(TimeFormat)
	var record TokenRecord
	record.Token = token
	record.UserName = userName
	record.Type = userType
	record.Created = created
	record.Expired = expired
	data, _ := json.Marshal(record)
	err := s.ioctx.SetOmap(TokenTable, map[string][]byte{
		token:data,
	})
	if err != nil {
		helper.Logger.Errorln("Error InsertTokenRecord", token, userName, userType, created, expired, err.Error())
	}
	return err
}

func (s *CephStore)RemoveTokenRecord(token string) error {
	err := s.ioctx.RmOmapKeys(TokenTable,[]string{token})
	if err != nil {
		helper.Logger.Errorln("Error remove token", token, err.Error())
	}
	return err
}

func (s *CephStore)SearchExistedToken(userName string) (TokenRecord, error) {
	var record TokenRecord
	result, err := s.ioctx.GetOmapValues(TokenTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Tokens", err)
		return record, err
	}
	for _, v := range result {
		err = json.Unmarshal(v, &record)
		if err != nil {
			return record, err
		}
		if record.UserName == userName {
			expired := time.Now().Add(time.Duration(helper.Config.TokenExpire * 1000000000)).Format(TimeFormat)
			record.Expired = expired
			data, _ := json.Marshal(record)
			err := s.ioctx.SetOmap(TokenTable, map[string][]byte{
				record.Token:data,
			})
			if err != nil {
				helper.Logger.Errorln("Error Update TokenRecord", record.Token, err.Error())
				return record, err
			}
			return record, nil
		}
	}
	return record, errors.New("token not found")
}

func (s *CephStore)ListExpiredTokens() ([]TokenRecord, error) {
	var records []TokenRecord = make([]TokenRecord, 0)
	result, err := s.ioctx.GetOmapValues(TokenTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Tokens", err)
		return records, err
	}
	now := time.Now()
	for _, v := range result {
		var record TokenRecord
		err = json.Unmarshal(v, &record)
		if err != nil {
			return nil, err
		}
		t, err := time.Parse(TimeFormat, record.Expired)
		if err != nil {
			helper.Logger.Errorln("Error parse expired time: ", err)
			records = append(records, record)
			continue
		}
		if now.Sub(t) > 0 {
			records = append(records, record)
		}

	}
	return records, err
}

func (s *CephStore)GetTokenRecord(token string) (TokenRecord, error) {
	var record TokenRecord
	results, err := s.ioctx.GetOmapValuesByKeys(TokenTable,[]string{token})
	if err != nil {
		return record, err
	}
	if len(results) > 0 {
		err = json.Unmarshal(results[token], &record)
		return record, err
	}
	return record, ErrTokenInvalid
}

func (s *CephStore)InsertSnmpRecord(snmpAddress string) error {
	var record SnmpRecord
	record.SnmpAddress = snmpAddress
	data, _ := json.Marshal(record)
	err := s.ioctx.SetOmap(SnmpTable, map[string][]byte{
		"snmp":data,
	})

	if err != nil {
		helper.Logger.Errorln("Error add snmp address", snmpAddress, err.Error())
		return err
	}
	return nil
}


func (s *CephStore)GetSnmpRecord() (SnmpRecord, error) {
	var record SnmpRecord
	results, err := s.ioctx.GetOmapValuesByKeys(SnmpTable,[]string{"snmp"})
	if err != nil {
		return record, err
	}
	if len(results) > 0 {
		err = json.Unmarshal(results["snmp"], &record)
		return record, err
	}
	return record, nil
}

func (s *CephStore)InsertNtpRecord(ntpAddress string) error {
	var record NtpRecord
	record.NtpAddress = ntpAddress
	data, _ := json.Marshal(record)
	err := s.ioctx.SetOmap(NtpTable, map[string][]byte{
		"ntp":data,
	})

	if err != nil {
		helper.Logger.Errorln("Error add ntp address", ntpAddress, err.Error())
		return err
	}
	return err
}


func (s *CephStore)GetNtpRecord() (NtpRecord, error) {
	var record NtpRecord
	results, err := s.ioctx.GetOmapValuesByKeys(NtpTable,[]string{"ntp"})
	if err != nil {
		return record, err
	}
	if len(results) > 0 {
		err = json.Unmarshal(results["ntp"], &record)
		return record, err
	}
	return record, nil
}

func (s *CephStore)GetHistoryRecord() ([]HistoryRecord, error) {
	var records []HistoryRecord = make([]HistoryRecord, 0)
	result, err := s.ioctx.GetOmapValues(HistoryTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Tokens", err)
		return records, err
	}
	for _, v := range result {
		var record HistoryRecord
		err = json.Unmarshal(v, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, err
}

func (s *CephStore)InsertOpRecord(userName, opType string) error {
	var record HistoryRecord
	created := time.Now().Format(TimeFormat)
	record.UserName = userName
	record.OpType = opType
	record.Time = created
	data, _ := json.Marshal(record)
	err := s.ioctx.SetOmap(HistoryTable, map[string][]byte{
		created:data,
	})
	if err != nil {
		helper.Logger.Errorln("Error log operate record", userName, opType, err.Error())
	}
	return err
}

func hexTimeStringToDate(s string) string {
	t, _ := strconv.ParseInt(s, 16, 64)
	return time.Unix(t,0).Format(TimeFormat)
}

//golang get datetime type from mysql
//https://www.jianshu.com/p/444bf0fddcd7
//https://github.com/go-sql-driver/mysql#timetime-support
//I use string for mysql.DATETIME
func (s *CephStore)GetAlertRecord() ([]AlertRecord, error) {
	var records []AlertRecord = make([]AlertRecord, 0)
	result, err := s.ioctx.GetOmapValues(SendedTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Alert", err)
		return records, err
	}
	for k, v := range result {
		var record AlertRecord
		splited := strings.Split(k, ":")
		if len(splited) != 3 {
			helper.Logger.Errorln("Error Split alert key, bad key struct :", k)
			continue
		}
		record.EventName = splited[0]
		record.Node = splited[1]
		record.SendTime = splited[2]
		record.EventId = string(v[:])
		record.Status = "sended"
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		first, _ := strconv.ParseInt(records[i].SendTime, 16, 64)
		second, _ := strconv.ParseInt(records[j].SendTime, 16, 64)
		return first > second
	})
	for i:=0; i < len(records); i++ {
		records[i].SendTime = hexTimeStringToDate(records[i].SendTime)
	}
	//rows, err := s.db.Query("select eventid, eventname, node, sendtime, status  from alerts order by sendtime DESC limit 100")
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var alert AlertRecord
	//	rows.Scan(&alert.EventId, &alert.EventName, &alert.Node, &alert.SendTime, &alert.Status)
	//	alerts = append(alerts, alert)
	//}
	return records, nil

}

func (s *CephStore)GetResolvedAlertRecord() ([]AlertRecord, error) {
	var records []AlertRecord = make([]AlertRecord, 0)
	result, err := s.ioctx.GetOmapValues(ResolvedTable, "","", math.MaxUint32)
	if err != nil {
		helper.Logger.Errorln("Error List Alert", err)
		return records, err
	}
	for k, v := range result {
		var record AlertRecord
		splited := strings.Split(k, ":")
		if len(splited) != 3 {
			helper.Logger.Errorln("Error Split alert key, bad key struct :", k)
			continue
		}
		record.EventName = splited[0]
		record.Node = splited[1]
		record.SendTime = splited[2]
		record.EventId = string(v[:])
		record.Status = "resolved"
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		first, _ := strconv.ParseInt(records[i].SendTime, 16, 64)
		second, _ := strconv.ParseInt(records[j].SendTime, 16, 64)
		return first > second
	})
	for i:=0; i < len(records); i++ {
		records[i].SendTime = hexTimeStringToDate(records[i].SendTime)
	}
	//rows, err := s.db.Query("select eventid, eventname, node, sendtime, status  from alerts order by sendtime DESC limit 100")
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var alert AlertRecord
	//	rows.Scan(&alert.EventId, &alert.EventName, &alert.Node, &alert.SendTime, &alert.Status)
	//	alerts = append(alerts, alert)
	//}
	return records, nil

}

func (s *CephStore)Close() {
	s.conn.Shutdown()
}