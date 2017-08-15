package api

import (
"gopkg.in/iris.v4"
"github.com/journeymidnight/yig-iam/helper"
"github.com/journeymidnight/yig-iam/db"
. "github.com/journeymidnight/yig-iam/api/datatype"
)

func CreateProject(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	id := "p-" + string(helper.GenerateRandomId())
	err := db.InsertProjectRecord(id, query.ProjectName, tokenRecord.AccountId, query.Description)
	if err != nil {
		helper.Logger.Println(5, "failed CreateProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed CreateProject",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DeleteProject(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	err := db.RemoveProjectRecord(query.ProjectId, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DeleteProject",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DescribeProject(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	record, err := db.DescribeProjectRecord(query.ProjectId, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeProject",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:record})
	return
}

func ListProjects(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListProjects, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	if tokenRecord.Type == ROLE_ACCOUNT || tokenRecord.Type == ROLE_ROOT {
		records, err := db.ListProjectRecords(tokenRecord.AccountId)
		if err != nil {
			helper.Logger.Println(5, "failed DescribeProject for query:", query)
			c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DescribeProject",Data:query})
			return
		}

		var resp ListProjectResp
		resp.Projects = records
		resp.Limit = 20
		resp.Offset = 0
		resp.Total = len(records)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:resp})
	} else {

	}
	return
}
