package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	//	"time"
	//	"syscall"
	"gopkg.in/iris.v4"
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/log"
	"github.com/journeymidnight/yig-iam/api"
	"github.com/journeymidnight/yig-iam/db"
	tokenMiddleware "github.com/journeymidnight/yig-iam/middleware/token"
	"github.com/hsluoyz/casbin"
	"github.com/casbin/xorm-adapter"
	_ "github.com/go-sql-driver/mysql"
)

var logger *log.Logger

func main() {
	fmt.Println(5, "enter 2:")
	helper.SetupConfig()
	fmt.Println(5, "enter 3:")
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	fmt.Println(6, "enter 3:", helper.CONFIG.PidFile)
	if err := helper.CreatePidfile(helper.CONFIG.PidFile); err != nil {
		fmt.Printf("can not create pid file %s\n", helper.CONFIG.PidFile)
		return
	}
	fmt.Println(5, "enter 4:")
	defer helper.RemovePidfile(helper.CONFIG.PidFile)

	/* log  */
	f, err := os.OpenFile(helper.CONFIG.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file " + helper.CONFIG.LogPath)
	}
	defer f.Close()
	fmt.Println(5, "enter 5:")
	logger = log.New(f, "[yig]", log.LstdFlags, helper.CONFIG.LogLevel)
	helper.Logger = logger
	fmt.Println(5, "enter 0:")

	a := xormadapter.NewAdapter("mysql", "root:12345678@tcp(127.0.0.1:3306)/")
	helper.Enforcer = casbin.NewEnforcer("./config/basic_model.conf", a)
	helper.Enforcer.LoadPolicy()
	roles := helper.Enforcer.GetAllRoles()
	fmt.Println(5, "enter 1:", len(roles))
	if len(roles) == 0 {
		logger.Println(5, "roles number:", len(roles))
		helper.Casbin_init()
	}
	/* redirect stdout stderr to log  */
	syscall.Dup2(int(f.Fd()), 2)
	syscall.Dup2(int(f.Fd()), 1)
	db.Db = db.CreateDbConnection()
	defer db.Db.Close()
	tokenMiddleware := tokenMiddleware.New()
	iris.Post("/", tokenMiddleware.Serve, api.ApiHandle)
	iris.Listen(":"+strconv.Itoa(helper.CONFIG.BindPort))
}

