package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	. "github.com/journeymidnight/nier/src/api/datatype"
	. "github.com/journeymidnight/nier/src/error"
	"github.com/journeymidnight/nier/src/helper"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const CreateViewBinary = "cd /opt/cephdeploy && /usr/bin/fab AddOneExporter:dirname='%s'"
const RemoveViewBinary = "cd /opt/cephdeploy && /usr/bin/fab RemoveOneExporter:dirname='%s'"
const exportsFSFile = "/etc/exports"

func parseCurrentExportsName() ([]string, error) {
	f, err := os.Open(exportsFSFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var exportedFS []string = make([]string,0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		stringList := strings.Split(scanner.Text(), " ")
		if len(stringList) > 1 {
			viewName := stringList[0][1:]
			exportedFS = append(exportedFS, viewName)
		}
	}

	return exportedFS, nil
}

func parseCurrentExportsPath() ([]string, error) {
	f, err := os.Open(exportsFSFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var exportedFS []string = make([]string,0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		stringList := strings.Split(scanner.Text(), " ")
		if len(stringList) > 1 {
			exportedFS = append(exportedFS, stringList[0])
		}
	}

	return exportedFS, nil
}

func ListView(w http.ResponseWriter, r *http.Request) {
	exportedFS, err := parseCurrentExportsPath()
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(exportedFS))
}

func CreateView(w http.ResponseWriter, r *http.Request) {

	//TODO: Need A LOCK
	body, _ := ioutil.ReadAll(r.Body)
	createView := CreateViewRecords{}
	err := json.Unmarshal(body, &createView)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	//check: not to create the same exported directory
	exportedFS, err := parseCurrentExportsName()
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}

	for _, v := range exportedFS {
		if v == createView.ExportedFS {
			WriteErrorResponse(w, r, ErrDuplicatedView)
			return
		}
	}

	command := fmt.Sprintf(CreateViewBinary, createView.ExportedFS)
	helper.Logger.Infof("Run command %s\n", command)
	out, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		helper.Logger.Errorf("Create view failed: result %s, %v", out, err)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Logger.Infof("Create view: ok %s, %v", out, err)
	WriteSuccessResponse(w, nil)
}

func RemoveView(w http.ResponseWriter, r *http.Request) {

	//TODO: Need A LOCK
	body, _ := ioutil.ReadAll(r.Body)
	removeView := RemoveViewRecord{}
	err := json.Unmarshal(body, &removeView)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	//check: not to create the same exported directory
	exportedFS, err := parseCurrentExportsName()
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}

	found := false
	for _, v := range exportedFS {
		if v == removeView.ExportedFS {
			found = true
			break
		}
	}

	if found == false {
		WriteErrorResponse(w, r, ErrNoneExistedView)
		return
	}

	command := fmt.Sprintf(RemoveViewBinary, removeView.ExportedFS)
	helper.Logger.Infof("Run command %s\n", command)
	out, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		helper.Logger.Errorf("Remove view failed: result %s, %v", out, err)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Logger.Infof("Remove view: ok %s, %v", out, err)
	WriteSuccessResponse(w, nil)
}
