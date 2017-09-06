package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	//	"time"
	//	"syscall"
	"github.com/casbin/xorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hsluoyz/casbin"
	"github.com/journeymidnight/yig-iam/api"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/log"
	tokenMiddleware "github.com/journeymidnight/yig-iam/middleware/token"
	"gopkg.in/iris-contrib/middleware.v4/cors"
	//"gopkg.in/kataras/iris.v4"
	"gopkg.in/kataras/iris.v4"
)

var logger *log.Logger

func main() {
	fmt.Println(5, "enter 2:")
	helper.SetupConfig()
	fmt.Println(5, "enter 3:")
	defer func() {
		if err := recover(); err != nil {
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

	a := xormadapter.NewAdapter("mysql", helper.CONFIG.CasbinDbString)
	helper.Enforcer = casbin.NewEnforcer("./config/basic_model.conf", a)
	helper.Enforcer.LoadPolicy()
	roles := helper.Enforcer.GetAllRoles()
	fmt.Println(5, "enter 1:", len(roles))
	//if len(roles) == 0 {
	//	logger.Println(5, "roles number:", len(roles))
	helper.Casbin_init()
	//}
	/* redirect stdout stderr to log  */
	syscall.Dup2(int(f.Fd()), 2)
	syscall.Dup2(int(f.Fd()), 1)
	db.Db = db.CreateDbConnection()
	defer db.Db.Close()
	tokenMiddleware := tokenMiddleware.New()

	app := iris.New()
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"Content-Type", "X-Iam-Token", "X-Le-Endpoint", "X-Le-Key", "X-Le-Secret"},
	})
	app.Use(c)
	app.Post("/iamapi", tokenMiddleware.Serve, api.ApiHandle)
	app.Post("/losapi", tokenMiddleware.Serve, api.LosApiHandler)
	app.Get("/env", api.EnvHandler)
	app.Listen(":" + strconv.Itoa(helper.CONFIG.BindPort))
}
