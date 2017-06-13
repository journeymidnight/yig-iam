package api

import (
	"gopkg.in/iris.v4"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/db"
)

func CreateAccount(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	helper.Logger.Println(5, "CreateAccount", tokenRecord.UserName)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_CreateAccount, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	var accountId []byte
	loop := 0
	for loop < 3 {
		helper.Logger.Println(5, "CreateAccount ENTER LOOP")
		accountId = helper.GenerateRandomNumberId()
		helper.Logger.Println(5, "CreateAccount ENTER LOOP", string(accountId))
		exist, err := db.CheckAccountIdExist(string(accountId))
		helper.Logger.Println(5, "CreateAccount ENTER LOOP1 ", exist, err)
		if exist == false && err == nil{
			break
		}
		loop = loop + 1
	}
	helper.Logger.Println(5, "CreateAccount OUT LOOP", loop)
	if loop >= 3 {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed create account",Data:query})
		return
	}

	err := db.InsertUserRecord(query.UserName, query.Password, ROLE_ACCOUNT, query.Email, query.DisplayName, string(accountId))
	if err != nil {
		helper.Logger.Println(5, "failed create account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed create account:" + err.Error(),Data:query})
		return
	}
	helper.Enforcer.AddRoleForUser(query.UserName, ROLE_ACCOUNT)
	helper.Enforcer.SavePolicy()
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DeleteAccount(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DeleteAccount, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	record, err := db.DescribeAccount(query.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed search account",Data:query})
		return
	}

	err = db.RemoveAccountId(query.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed delete account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed delete account",Data:query})
		return
	}
	helper.Logger.Println(5, "DeleteAccount:", query.UserName)
	helper.Enforcer.DeleteRolesForUser(record.UserName)
	helper.Enforcer.SavePolicy()
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:""})
	return
}

func DescribeAccount(c *iris.Context, query QueryRequest) {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeAccount, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	record, err := db.DescribeAccount(query.AccountId)
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed search account",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:record})
	return
}

func ListAccounts(c *iris.Context, query QueryRequest)  {
	tokenRecord := c.Get("token").(TokenRecord)
	if helper.Enforcer.Enforce(tokenRecord.UserName, API_DescribeAccount, ACT_ACCESS) != true {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4030,Message:"You do not have permission to perform", Data:query})
		return
	}
	records, err := db.ListAccountRecords()
	if err != nil {
		helper.Logger.Println(5, "failed search account for query:", query)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4010,Message:"failed search account",Data:query})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:records})
	return
}
