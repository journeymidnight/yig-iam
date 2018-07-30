package api

import (
	"encoding/json"
	. "github.com/journeymidnight/nier/src/api/datatype"
	"github.com/journeymidnight/nier/src/store"
	. "github.com/journeymidnight/nier/src/error"
	"github.com/journeymidnight/nier/src/helper"
	"database/sql"
	"io/ioutil"
	"net/http"
	"github.com/Jeffail/gabs"
)

const (
	DEPLOY_CONFIG_PATH = "/opt/cephdeploy/config.json"
)

func SetSNMP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	err = store.Db.InsertSnmpRecord(req.SnmpAddress)
	if err != nil {
		helper.Logger.Infoln("failed SetSNMP for query:", req)
		WriteErrorResponse(w, r, err)
		return
	}
	token, ok := r.Context().Value(REQUEST_TOKEN_KEY).(TokenRecord)
	if ok != false {
		store.Db.InsertOpRecord(token.UserName, OP_TYPE_SET_SNMP)
	}
	WriteSuccessResponse(w, nil)
}

func GetSNMP(w http.ResponseWriter, r *http.Request) {
	record, err := store.Db.GetSnmpRecord()
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Infoln("failed GetSNMP for query:", r)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(record))
}

func SetNTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequestNtp{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	err = store.Db.InsertNtpRecord(req.NtpAddress)
	if err != nil {
		helper.Logger.Infoln("failed SetNTP for query:", req)
		WriteErrorResponse(w, r, err)
		return
	}
	token, ok := r.Context().Value(REQUEST_TOKEN_KEY).(TokenRecord)
	if ok != false {
		store.Db.InsertOpRecord(token.UserName, OP_TYPE_SET_NTP)
	}
	WriteSuccessResponse(w, nil)
}

func GetNTP(w http.ResponseWriter, r *http.Request) {
	record, err := store.Db.GetNtpRecord()
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Infoln("failed GetNTP for query:", r)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(record))
}

func ListAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := store.Db.GetAlertRecord()
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Infoln("failed ListAlert for query:", r)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(alerts))
}

func ListResolvedAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := store.Db.GetResolvedAlertRecord()
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Infoln("failed ListAlert for query:", r)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(alerts))
}


func GetVip(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile(DEPLOY_CONFIG_PATH)
	if err != nil {
		WriteErrorResponse(w, r, err)
	}
	jsonParsed, err := gabs.ParseJSON(content)
	var value string
	var ok bool
	value, ok = jsonParsed.Path("vip").Data().(string)
	if ok {
		WriteSuccessResponse(w, EncodeResponse(VipRecord{value}))
	} else {
		WriteErrorResponse(w, r, ErrFailedGetVip)
	}
}
