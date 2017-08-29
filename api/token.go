package api

import (
	"github.com/google/uuid"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"gopkg.in/iris.v4"
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
				c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "user name or password incorrect", Data: ""})
				return
			}
		}
	}

	if record.Status == "inactive" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "your account has been disabled by administrator", Data: ""})
		return
	}

	uuid := uuid.New()
	helper.Logger.Println(5, "ConnectService uuid length:", len(uuid.String()))
	err = db.InsertTokenRecord(uuid.String(), record.UserName, record.AccountId, record.Type)
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "InsertTokenRecord error", Data: ""})
		return
	}
	var resp ConnectServiceResponse
	resp.Token = uuid.String()
	resp.Type = record.Type
	resp.AccountId = record.AccountId

	//for accounts, return linked projects
	if resp.Type != ROLE_ROOT {
		prjs, err := db.GetLinkedProjects(record.AccountId)
		if err != nil || len(prjs) < 1 {
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4010, Message: "cann't find linked project to this user", Data: ""})
			return
		}

		var found bool = false
		if query.ProjectId != "" {
			for _, p := range prjs {
				if p.ProjectId == query.ProjectId {
					found = true
					break
				}
			}
		}

		if found == true {
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
			return
		} else {
			resp.ProjectSet = prjs
			c.JSON(iris.StatusOK, QueryResponse{RetCode: 4102, Message: "please choose a project", Data: resp})
			return
		}
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: resp})
	return
}
