package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	//	"time"
	//	"syscall"
	"github.com/kataras/iris"
	"legitlab.letv.cn/yig/iam/helper"
	"legitlab.letv.cn/yig/iam/log"
	"legitlab.letv.cn/yig/iam/api"
	"legitlab.letv.cn/yig/iam/db"
	tokenMiddleware "legitlab.letv.cn/yig/iam/middleware/token"
	"github.com/hsluoyz/casbin"
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
	helper.Enforcer = casbin.NewEnforcer("./config/casbin.conf")
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

