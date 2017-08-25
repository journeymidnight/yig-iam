package api

import (
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/iris.v4"
)

func CreateProjectRole(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateProjectRole, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	err := db.InsertProjectRoleRecord(query.ProjectId, query.AccountId, query.Role)
	if err != nil {
		helper.Logger.Println(5, "failed CreateProjectRole for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DeleteProjectRole(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteProjectRole, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	err := db.RemoveProjectRoleRecord(query.ProjectId, query.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DeleteProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DescribeProjectRoles(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeProjectRoles, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	records, err := db.ListProjectRoleRecordsByProjectId(query.ProjectId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProjectRoles for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProjectRoles", Data: query})
		return
	}

	// got account info
	//for _, pr := range records {
	for i := 0; i < len(records); i++ {
		ur, _ := db.DescribeAccount(records[i].UserId)
		records[i].Email = ur.Email
		records[i].Username = ur.UserName
	}

	var resp ListProjectRoleResp
	resp.RoleSet = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}

//func ListProjects(c *iris.Context, query QueryRequest) {
//	tokenRecord := c.Get("token").(TokenRecord)
//	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListProjects, ACT_ACCESS) != true {
//		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
//		return
//	}
//	if tokenRecord.Type == ROLE_ACCOUNT || tokenRecord.Type == ROLE_ROOT {
//		records, err := db.ListProjectRecords(tokenRecord.AccountId)
//		if err != nil {
//			helper.Logger.Println(5, "failed DescribeProject for query:", query)
//			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
//			return
//		}
//
//		var resp ListProjectResp
//		resp.Projects = records
//		resp.Limit = 20
//		resp.Offset = 0
//		resp.Total = len(records)
//		c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
//	} else {
//
//	}
//	return
//}
//
