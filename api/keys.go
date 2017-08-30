package api

import (
	"database/sql"
	"errors"

	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/iris.v4"
)

func createkeypair(ProjectId, AccountId, KeyName, Description string, IsAutoGen bool) (string, string, error) {
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
		return "", "", errors.New("failed to find a properiate ak")
	}
	accessSecret := helper.GenerateRandomIdByLength(40)

	ak := string(accessKey[:])
	sk := string(accessSecret[:])
	err := db.InsertAkSkRecord(ak, sk, ProjectId, AccountId, KeyName, Description, IsAutoGen)
	if err != nil {
		helper.Logger.Println(5, "failed CreateAccessKey for project :", ProjectId)
		return "", "", errors.New("failed CreateAccessKey, maybe you create two keys with same name")
	}

	return ak, sk, nil
}

func CreateAccessKey(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateAccessKey, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}

	ak, sk, err := createkeypair(query.ProjectId, tokenRecord.AccountId, query.KeyName, query.Description, false)
	if err != nil {
		helper.Logger.Printf(5, "failed CreateAccessKey for query %+v, error is %s", query, err.Error())
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: err.Error(), Data: query})
		return
	}

	// create a key in s3
	err = s3CreateKey(query.ProjectId, ak, sk) // fixme: use username as project id for now
	if err != nil {
		helper.Logger.Println(5, "failed CreateAccessKey in s3")
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateAccessKey, failed to connect to object storage system", Data: query})
		return
	}

	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DeleteAccessKey(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteAccessKey, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}

	// delete it from s3 before delete it
	// we must delete it from s3 before delete from db, otherwise we will never delete the key successfully because we can't got the sk and pid anymore.
	sk, pid, err := db.GetSkAndProjectByAk(query.AccessKey)
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "no such pair of key", Data: ""})
		return
	} else {
		// delete the in s3
		err = s3DeleteKey(pid, query.AccessKey, sk)
		if err != nil {
			helper.Logger.Println(5, "failed s3DeleteKey in s3")
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateAccessKey, failed to connect to object storage system", Data: query})
			return
		}
	}

	err = db.RemoveAkSkRecord(query.AccessKey)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteAccessKey for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DeleteAccessKey", Data: query})
		return
	}

	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func ListAccessKeysByProject(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListAccessKeysByProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	records, err := db.ListAkSkRecordByProject(query.ProjectId)
	if err != nil {
		helper.Logger.Println(5, "failed ListAccessKeysByProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed ListAccessKeysByProject", Data: query})
		return
	}

	var resp DescribeKeysResp
	resp.AccessKeySet = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)

	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}

func DescribeAccessKeys(c *iris.Context, query QueryRequest) {
	if c.RequestHeader("X-Le-Key") != helper.CONFIG.ManageKey || c.RequestHeader("X-Le-Secret") != helper.CONFIG.ManageSecret {
		helper.Logger.Println(5, "unauthorized request")
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "unauthorized request", Data: ""})
		return
	}
	var resp DescribeKeysResp
	records, err := db.GetKeysByAccessKeys(query.AccessKeys)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
		return
	}
	resp.AccessKeySet = records
	resp.Limit = len(records)
	resp.Offset = 0
	resp.Total = resp.Limit
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}

func DescribeAccessKeysWithToken(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, ACTION_DescribeAccessKeysWithToken, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	var resp DescribeKeysResp
	records, err := db.GetKeysByAccount(tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
		return
	}
	resp.AccessKeySet = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}

func GetAutogenkeysByProjectId(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListAccessKeysByProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	record, err := db.GetAutogenAkSkRecordByProject(query.ProjectId)
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Println(5, "failed GetAutogenkeysByProjectId for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed GetAutogenkeysByProjectId ", Data: query})
		return
	} else if err == sql.ErrNoRows {
		ak, sk, err := createkeypair(query.ProjectId, tokenRecord.AccountId, "autogen", "autogen", true)
		if err != nil {
			helper.Logger.Printf(5, "failed GetAutogenkeysByProjectId for query %+v, error is %s", query, err.Error())
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: err.Error(), Data: query})
			return
		}

		// create a key in s3
		err = s3CreateKey(query.ProjectId, ak, sk)
		if err != nil {
			helper.Logger.Println(5, "failed CreateAccessKey in s3")
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateAccessKey, failed to connect to object storage system", Data: query})
			return
		}

		record.AccessKey = ak
		record.AccessSecret = sk
	}

	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: record})
	return
}
