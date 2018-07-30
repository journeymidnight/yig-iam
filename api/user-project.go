package api

import (
	. "github.com/journeymidnight/yig-iam/api/datatype"
	. "github.com/journeymidnight/yig-iam/error"
	"github.com/journeymidnight/yig-iam/db"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func LinkUserWithProject(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	err = db.LinkUserWithProject(query.ProjectId, query.UserId, query.Acl)
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

func UnLinkUserWithProject(w http.ResponseWriter, r *http.Request)  {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	err = db.UnlinkUserWithProject(query.ProjectId, query.UserId)
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, nil)
	return
}

//func ListProjectByUser(w http.ResponseWriter, r *http.Request)  {
//	tokenRecord := c.Get("token").(TokenRecord)
//	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListProjectByUser, ACT_ACCESS) != true {
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
//		return
//	}
//	var realName string
//	if tokenRecord.Type == ROLE_USER {
//		realName = tokenRecord.UserName
//	} else {
//		realName = tokenRecord.AccountId + ":" + query.UserName
//	}
//	record, err := db.ListProjects(realName)
//	if err != nil {
//		helper.Logger.Println(5, "failed search account for query:", query)
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed search account",Data:query})
//		return
//	}
//	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:record})
//	return
//}

func ListUserByProject(w http.ResponseWriter, r *http.Request)  {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	record, err := db.ListUsersByProject(query.ProjectId)
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(record))
	return
}
