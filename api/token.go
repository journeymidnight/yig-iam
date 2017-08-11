package api

import (
	"gopkg.in/iris.v4"
	"github.com/google/uuid"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
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
		//retry auth by email
		if query.Email != "" {
			helper.Logger.Println(5, "ConnectService user email:", query.Email, query.Password)
			record, err = db.ValidEmailAndPassword(query.Email, query.Password)
			if err != nil {
				c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"user name or password incorrect",Data:""})
				return
			}
		}
	}

	if record.Status == "inactive" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"your account has been disabled by administrator",Data:""})
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
