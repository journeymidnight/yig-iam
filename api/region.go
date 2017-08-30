package api

import (
	"fmt"

	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/iris.v4"
)

func CreateRegion(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateRegion, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	id := "r-" + string(helper.GenerateRandomId())
	err := db.InsertRegionRecord(id, query.RegionName)
	if err != nil {
		helper.Logger.Println(5, "failed CreateProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed CreateProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func ModifyRegionAttributes(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_ModifyRegionAttributes, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}

	err := db.UpdateRegionRecord(query.RegionId, query.RegionName)
	fmt.Println(query.RegionId, query.RegionName)
	if err != nil {
		helper.Logger.Printf(5, "failed modify project for query: %+v", query, "regionid is %s, regionname is %s\r\n", query.RegionId, query.RegionName)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed modify project", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DeleteRegion(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteRegion, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	err := db.RemoveRegionRecord(query.RegionId)
	if err != nil {
		helper.Logger.Println(5, "failed DeleteRegion for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DeleteRegion", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: ""})
	return
}

func DescribeRegion(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeProject, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}
	record, err := db.DescribeProjectRecord(query.ProjectId, tokenRecord.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: record})
	return
}

func DescribeRegions(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeRegions, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
		return
	}

	records, err := db.ListRegionRecords()
	if err != nil {
		helper.Logger.Println(5, "failed DescribeProject for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "failed DescribeProject", Data: query})
		return
	}

	var resp ListRegionResp
	resp.Regions = records
	resp.Limit = 20
	resp.Offset = 0
	resp.Total = len(records)
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}
