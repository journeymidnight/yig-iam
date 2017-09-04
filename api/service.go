package api

import (
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/iris.v4"
)

func CreateService(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateService, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	id := "s-" + string(helper.GenerateRandomId())
	err := db.InsertServiceRecord(id, query.PublicUrl, query.Endpoint, query.RegionId)
	if err != nil {
		helper.Logger.Println(5, "failed CreateProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func ModifyServiceAttributes(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ModifyServiceAttributes, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	err := db.UpdateServiceRecord(query.ServiceId, query.Endpoint)
	if err != nil {
		helper.Logger.Printf(5, "failed modify service for query: %+v", query, "regionid is %s, regionname is %s\r\n", query.ServiceId, query.Endpoint)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed modify project", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DeleteService(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteService, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	err := db.RemoveServiceRecord(query.ServiceId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteRegion for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DeleteRegion", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DescribeServices(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeServices, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}

	records, err := db.ListSerivceRecords()
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
		return
	}

	var resp ListServiceResp
	resp.Services = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})

	return
}
