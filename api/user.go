/*
Table Account
accountName
password
description
type
status
created
updated

Table User
userName
password
email
accountName
Projects
status
created
updated

Table Project
projectID
accountName
services
quota
status
created
updated

Table Keys
key
secret
projectID
created
status

Table Token
uid
token
expired
created
 */
package api

import (
	"github.com/kataras/iris"
	"legitlab.letv.cn/yig/iam/helper"
	"legitlab.letv.cn/yig/iam/db"
	. "legitlab.letv.cn/yig/iam/api/datatype"
)

func CreateUser(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateUser, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	realName := tokenRecord.AccountId + ":" + query.UserName
	err := db.InsertUserRecord(realName, query.Password, ROLE_USER, query.Email, query.DisplayName, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed CreateUser for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateUser",Data:query})
		return
	}
	helper.Enforcer.AddRoleForUser(realName, ROLE_USER)
	helper.Enforcer.SavePolicy()
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DeleteUser(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteUser, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	realUserName := tokenRecord.AccountId + ":" + query.UserName
	err := db.RemoveUserRecord(realUserName, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteUser for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DeleteUser",Data:query})
		return
	}
	helper.Enforcer.DeleteRolesForUser(realUserName)
	helper.Enforcer.SavePolicy()
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DescribeUser(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeUser, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	var err error
	var record UserRecord
	if tokenRecord.Type == ROLE_ACCOUNT {
		record, err = db.DescribeUserRecord(tokenRecord.AccountId + ":" + query.UserName, tokenRecord.AccountId)
	} else {
		record, err = db.DescribeUserRecord(tokenRecord.UserName, tokenRecord.AccountId)
	}

	if err != nil {
		helper.Logger.Println(5, "failed DescribeUsert for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeUser",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:record})
	return
}

func ListUsers(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListUsers, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	records, err := db.ListUserRecords(tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed search account",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:records})
	return
}