package api

import (
	"github.com/kataras/iris"
	"legitlab.letv.cn/yig/iam/helper"
	. "legitlab.letv.cn/yig/iam/api/datatype"
	"legitlab.letv.cn/yig/iam/db"
)

func AddProjectService(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_AddProjectService, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	//TODO:valid project`s owner and service
	err := db.InsertProjectServiceRecord(query.ProjectId, query.Service, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failedAddProjectService:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed AddProjectService",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DelProjectService(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DelProjectService, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	err := db.RemoveProjectServiceRecord(query.ProjectId, query.Service, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DelProjectService for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed DelProjectService",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func ListServiceByProject(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ListServiceByProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	records, err := db.ListProjectServiceRecordByProject(query.ProjectId, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed ListServiceByProject:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed ListServiceByProject",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:records})
	return
}
