package api

import (
	"net/http"
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/db"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	. "github.com/journeymidnight/yig-iam/error"
	"io/ioutil"
	"encoding/json"
)

//func CreateAccessKey(w http.ResponseWriter, r *http.Request) {
//	tokenRecord := c.Get("token").(TokenRecord)
//	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateAccessKey, ACT_ACCESS) != true {
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
//		return
//	}
//	var accessKey []byte
//	i := 0
//	for i < 3 {
//		accessKey = helper.GenerateRandomIdByLength(20)
//		existed := db.IfAKExisted(string(accessKey[:]))
//		if existed == false {
//			break
//		}
//		i = i + 1
//	}
//	if i >= 3 {
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateAccessKey", Data:query})
//		return
//	}
//	accessSecret := helper.GenerateRandomIdByLength(40)
//
//	if query.ProjectId == "" {
//		query.ProjectId = "s3defaultproject"
//	}
//
//	err := db.InsertAkSkRecord(string(accessKey[:]), string(accessSecret[:]), query.ProjectId, tokenRecord.AccountId, query.KeyName, query.Description)
//	if err != nil {
//		helper.Logger.Println(5, "failed CreateAccessKey for query:", query)
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateAccessKey, maybe you create two keys with same name",Data:query})
//		return
//	}
//	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
//	return
//}

//func DeleteAccessKey(w http.ResponseWriter, r *http.Request) {
//	tokenRecord := c.Get("token").(TokenRecord)
//	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteAccessKey, ACT_ACCESS) != true {
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
//		return
//	}
//	err := db.RemoveAkSkRecord(query.AccessKey, tokenRecord.AccountId)
//	if err != nil {
//		helper.Logger.Println(5, "failed DeleteAccessKey for query:", query)
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DeleteAccessKey",Data:query})
//		return
//	}
//	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
//	return
//}

func FetchSecretKey(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if (query.AccessKey == "" && query.ProjectId == "") || (query.AccessKey != "" && query.ProjectId != ""){
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	if r.Header.Get("X-Le-Key") != helper.Config.ManageKey || r.Header.Get("X-Le-Secret") != helper.Config.ManageSecret {
		helper.Logger.Println(5, "unauthorized request")
		WriteErrorResponse(w, r, ErrNotAuthorised)
		return
	}
	resp := FetchAccessKeysResp{}

	if query.AccessKey != "" {
		items, err := db.GetKeyItemByAccessKey(query.AccessKey)
		if err != nil {
			WriteErrorResponse(w, r, err)
		} else {
			resp.AccessKeySet = items
			WriteSuccessResponse(w, EncodeResponse(resp))
		}
	} else {
		items, err := db.GetKeyItemsByProject(query.ProjectId)
		if err != nil {
			WriteErrorResponse(w, r, err)
		} else {
			resp.AccessKeySet = items
			WriteSuccessResponse(w, EncodeResponse(resp))
		}
	}
	return
}

//This is a legacy api for yig
func DescribeKeys(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if (query.AccessKeys == nil && query.ProjectId == "") || (query.AccessKeys != nil && query.ProjectId != ""){
		WriteErrorResponse(w, r, ErrInvalidParameters)
		return
	}
	if r.Header.Get("X-Le-Key") != helper.Config.ManageKey || r.Header.Get("X-Le-Secret") != helper.Config.ManageSecret {
		helper.Logger.Println(5, "unauthorized request")
		WriteErrorResponse(w, r, ErrNotAuthorised)
		return
	}
	resp := QueryResp{}

	if query.AccessKeys != nil {
		items, err := db.GetKeyItemsByAccessKeys(query.AccessKeys)
		if err != nil {
			WriteSuccessResponse(w, EncodeResponse(QueryRespAll{"failed DescribeAccessKeys", query, 4010}))
		} else {
			resp.AccessKeySet = items
			resp.Limit = len(items)
			resp.Offset = 0
			resp.Total = resp.Limit
			WriteSuccessResponse(w, EncodeResponse(QueryRespAll{"", resp, 0}))
		}
	} else {
		items, err := db.GetKeyItemsByProject(query.ProjectId)
		if err != nil {
			WriteSuccessResponse(w, EncodeResponse(QueryRespAll{"failed ListAccessKeysByProject", query, 4010}))
		} else {
			resp.AccessKeySet = items
			resp.Limit = len(items)
			resp.Offset = 0
			resp.Total = resp.Limit
			WriteSuccessResponse(w, EncodeResponse(QueryRespAll{"", resp, 0}))
		}
	}
	return
}
//
//func FetchAccessKeysWithToken(w http.ResponseWriter, r *http.Request) {
//	token, ok := r.Context().Value(REQUEST_TOKEN_KEY).(TokenRecord)
//	if ok != false {
//		store.Db.InsertOpRecord(token.UserName, OP_TYPE_SET_NTP)
//	}
//	body, _ := ioutil.ReadAll(r.Body)
//	query := &QueryRequest{}
//	err := json.Unmarshal(body, query)
//	if err != nil {
//		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
//		return
//	}
//	var resp DescribeKeysResp
//	records, err := db.GetKeysByAccount(tokenRecord.AccountId)
//	if err != nil {
//		helper.Logger.Println(5, "failed DescribeProject for query:", query)
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeProject",Data:query})
//		return
//	}
//	resp.AccessKeySet = records
//	resp.Limit = 20
//	resp.Offset = 0
//	resp.Total = len(records)
//	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:resp})
//	return
//}