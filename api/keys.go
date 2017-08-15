package api

import (
	"gopkg.in/iris.v4"
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/db"
	. "github.com/journeymidnight/yig-iam/api/datatype"
)

func CreateAccessKey(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateAccessKey, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	var accessKey []byte
	i := 0
	for i < 3 {
		accessKey = helper.GenerateRandomIdByLength(20)
		existed := db.IfAKExisted(string(accessKey[:]))
		if existed == false {
			break
		}
		i = i + 1
	}
	if i >= 3 {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateAccessKey", Data:query})
		return
	}
	accessSecret := helper.GenerateRandomIdByLength(40)

	if query.ProjectId == "" {
		query.ProjectId = "s3defaultproject"
	}

	err := db.InsertAkSkRecord(string(accessKey[:]), string(accessSecret[:]), query.ProjectId, tokenRecord.AccountId, query.KeyName, query.Description)
	if err != nil {
		helper.Logger.Println(5, "failed CreateAccessKey for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateAccessKey, maybe you create two keys with same name",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DeleteAccessKey(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteAccessKey, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	err := db.RemoveAkSkRecord(query.AccessKey, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteAccessKey for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DeleteAccessKey",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func ListAccessKeysByProject(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListAccessKeysByProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	record, err := db.ListAkSkRecordByProject(query.ProjectId, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed ListAccessKeysByProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed ListAccessKeysByProject",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:record})
	return
}

func DescribeAccessKeys(c *iris.Context, query QueryRequest) {
	if c.RequestHeader("X-Le-Key") != helper.CONFIG.ManageKey || c.RequestHeader("X-Le-Secret") != helper.CONFIG.ManageSecret {
		helper.Logger.Println(5, "unauthorized request")
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"unauthorized request",Data:""})
		return
	}
	var resp DescribeKeysResp
	records, err := db.GetKeysByAccessKeys(query.AccessKeys)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeProject",Data:query})
		return
	}
	resp.AccessKeySet = records
	resp.Limit = len(records)
	resp.Offset = 0
	resp.Total = resp.Limit
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:resp})
	return
}


func DescribeAccessKeysWithToken(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, ACTION_DescribeAccessKeysWithToken, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	var resp DescribeKeysResp
	records, err := db.GetKeysByAccount(tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeProject",Data:query})
		return
	}
	resp.AccessKeySet = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:resp})
	return
}