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

const (

)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.ProjectName == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	//ak, sk := helper.GenerateKey()
	err = db.CreateProject(query.ProjectName, PUBLIC_PROJECT, token.AccountId, query.Description)
	if err != nil {
		helper.Logger.Println(5, "failed CreateProject for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

func DeleteProject(w http.ResponseWriter, r *http.Request)  {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.ProjectId == "" {
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	err = db.RemoveProject(query.ProjectId, token.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteProject for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

func DescribeProject(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.ProjectId == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	record, err := db.DescribeProject(query.ProjectId, token.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(record))
	return
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	if token.Type == ROLE_ACCOUNT || token.Type == ROLE_USER{
		resp, err := db.ListProjects(token.UserId)
		if err != nil {
			WriteErrorResponse(w, r, err)
			return
		}
		WriteSuccessResponse(w, EncodeResponse(resp))
	} else {
		WriteSuccessResponse(w, EncodeResponse([]ListProjectResp{}))
	}
	return
}

func ListProjectByUser(w http.ResponseWriter, r *http.Request) {
	//token, _ := r.Context().Value(REQUEST_TOKEN_KEY).(Token)
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if query.UserId == ""{
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	record, err := db.ListProjects(query.UserId)
	if err != nil {
		helper.Logger.Println(5, "failed list projects for query:", query)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(record))
	return
}
