package api

import (
	"github.com/kataras/iris"
	"github.com/google/uuid"
	"legitlab.letv.cn/yig/iam/helper"
	. "legitlab.letv.cn/yig/iam/api/datatype"
	"legitlab.letv.cn/yig/iam/db"
)

func ConnectService(c *iris.Context, query QueryRequest) {
	var userName string
	if query.AccountId != "" {
		userName = query.AccountId + ":" + query.UserName
	} else {
		userName = query.UserName
	}
	helper.Logger.Println(5, "ConnectService user password:", userName, query.Password)
	record, err := db.ValidUserAndPassword(userName, query.Password)
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"user name or password incorrect",Data:""})
		return
	}
	uuid := uuid.New()
	helper.Logger.Println(5, "ConnectService uuid length:", len(uuid.String()))
	err = db.InsertTokenRecord(uuid.String(), record.UserName, record.AccountId, record.Type)
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"InsertTokenRecord error",Data:""})
		return
	}
	var resp ConnectServiceResponse
	resp.Token = uuid.String()
	resp.Type = record.Type
	resp.AccountId = record.AccountId
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:resp})
	return

}