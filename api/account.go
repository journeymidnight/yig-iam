package api

import (
	"net/http"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	. "github.com/journeymidnight/yig-iam/error"
	"io/ioutil"
	"encoding/json"
)

func CreateAccount(w http.ResponseWriter, r *http.Request)  {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	if query.UserName == "" || query.Password == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}

	var accountId []byte
	loop := 0
	for loop < 3 {
		accountId = helper.GenerateRandomNumberId()
		exist, err := db.CheckAccountIdExist(string(accountId))
		helper.Logger.Println(5, "CreateAccount ENTER LOOP1 ", exist, err)
		if exist == false && err == nil{
			break
		}
		loop = loop + 1
	}
	helper.Logger.Println(5, "CreateAccount OUT LOOP", loop)
	if loop >= 3 {
		WriteErrorResponse(w, r, ErrInternalError)
		return
	}

	err = db.CreateUser(query.UserName, query.Password, ROLE_ACCOUNT, query.Email, query.DisplayName, string(accountId))
	if err != nil {
		WriteErrorResponse(w, r, ErrInternalError)
		return
	}
	helper.Enforcer.AddRoleForUser(query.UserName, ROLE_ACCOUNT)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
	return
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	//token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	if query.UserName == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}

	err = db.DeleteAccount(query.UserName)
	if err != nil {
		helper.Logger.Errorln("failed delete account for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Logger.Infoln("DeleteAccount:", query.UserName)
	helper.Enforcer.DeleteRolesForUser(query.UserName)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
	return
}

func DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserName == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	err = db.DeactivateAccount(query.UserName)
	if err != nil {
		helper.Logger.Errorln("failed deactivate account: ", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

func ActivateAccount(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserName == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	err = db.ActivateAccount(query.UserName)
	if err != nil {
		helper.Logger.Errorln("failed deactivate account: ", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

func DescribeAccount(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserName == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	account, err := db.DescribeAccount(query.UserName)
	if err != nil {
		helper.Logger.Errorln("failed search account for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(account))
	return
}

func ListAccounts(w http.ResponseWriter, r *http.Request)  {
	accounts, err := db.ListAccounts()
	if err != nil {
		helper.Logger.Errorln("failed list accounts", err.Error())
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(accounts))
	return
}
