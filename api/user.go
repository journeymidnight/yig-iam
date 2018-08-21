package api

import (
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/db"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"io/ioutil"
	"net/http"
	"encoding/json"
	. "github.com/journeymidnight/yig-iam/error"
)

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserName == "" || query.Password == "" {
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	err = db.CreateUser(query.UserName, query.Password, ROLE_USER, query.Email, query.DisplayName, token.AccountId, true)
	if err != nil {
		helper.Logger.Println(5, "failed CreateUser for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Enforcer.AddRoleForUser(query.UserName, ROLE_USER)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request)  {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	if query.UserName == "" {
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}

	if query.UserName == "root"{
		WriteErrorResponse(w, r, ErrNotAuthorised)
		return
	}

	err = db.RemoveUser(query.UserName, token.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteUser for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Enforcer.DeleteRolesForUser(query.UserName)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
	return
}

func DescribeUser(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserName == "" {
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	user, err := db.DescribeUser(query.UserName, token.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeUsert for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(user))
	return
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	users, err := db.ListUsers(token.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed ListUsers for account:", token.AccountId)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(users))
	return
}
