package api

import (
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/kataras/iris.v4"
)

func LinkUserWithProject(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_LinkUserWithProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	realName := tokenRecord.AccountId + ":" + query.UserName
	err := db.InsertUserProjectRecord(query.ProjectId, realName)
	if err != nil {
		helper.Logger.Println(5, "failed create account for query:", query, err)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed LinkUserWithProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func UnLinkUserWithProject(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_UnLinkUserWithProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	realName := tokenRecord.AccountId + ":" + query.UserName
	err := db.RemoveUserProjectRecord(query.ProjectId, realName)
	if err != nil {
		helper.Logger.Println(5, "failed delete account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed RemoveUserProjectRecord", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func ListProjectByUser(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListProjectByUser, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	var realName string
	if tokenRecord.Type == ROLE_USER {
		realName = tokenRecord.UserName
	} else {
		realName = tokenRecord.AccountId + ":" + query.UserName
	}
	record, err := db.ListUserProjectRecordByUser(realName)
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed search account", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: record})
	return
}

func ListUserByProject(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListProjectByUser, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	record, err := db.ListUserProjectRecordByProject(query.ProjectId)
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed search account", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: record})
	return
}
