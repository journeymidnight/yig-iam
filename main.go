package main

import (
	"fmt"
	//	"time"
	//	"syscall"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api"
	"github.com/journeymidnight/yig-iam/db"
	"github.com/casbin/casbin"
	"github.com/casbin/xorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/journeymidnight/yig-iam/middleware/auth"
	"github.com/journeymidnight/yig-iam/middleware/cors"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"strings"
	"time"
	"net/http"
)

var AllRoutes []map[string]string

//func ApiHandle(c *iris.Context) {
//	query := c.Get("queryRequest").(QueryRequest)
//
//	switch query.Action {
//	case ACTION_ConnectService:
//		ConnectService(c, query)
//	case ACTION_CreateAccount:
//		CreateAccount(c, query)
//	case ACTION_DeleteAccount:
//		DeleteAccount(c, query)
//	case ACTION_DescribeAccount:
//		DescribeAccount(c, query)
//	case ACTION_DeactivateAccount:
//		DeactivateAccount(c, query)
//	case ACTION_ActivateAccount:
//		ActivateAccount(c, query)
//	case ACTION_ListAccounts:
//		ListAccounts(c, query)
//	case ACTION_ListUsers:
//		ListUsers(c, query)
//	case ACTION_DescribeUser:
//		DescribeUser(c, query)
//	case ACTION_CreateUser:
//		CreateUser(c, query)
//	case ACTION_DeleteUser:
//		DeleteUser(c, query)
//	case ACTION_DescribeProject:
//		DescribeProject(c, query)
//	case ACTION_CreateProject:
//		CreateProject(c, query)
//	case ACTION_DeleteProject:
//		DeleteProject(c, query)
//	case ACTION_ListProjects:
//		ListProjects(c, query)
//	case ACTION_LinkUserWithProject:
//		LinkUserWithProject(c, query)
//	case ACTION_UnLinkUserWithProject:
//		UnLinkUserWithProject(c, query)
//	case ACTION_ListProjectByUser:
//		ListProjectByUser(c, query)
//	case ACTION_ListUserByProject:
//		ListUserByProject(c, query)
//	case ACTION_AddProjectService:
//		AddProjectService(c, query)
//	case ACTION_DelProjectService:
//		DelProjectService(c, query)
//	case ACTION_ListServiceByProject:
//		ListServiceByProject(c, query)
//	case ACTION_DescribeAccessKeys:
//		DescribeAccessKeys(c, query)
//	case ACTION_DescribeAccessKeysWithToken:
//		DescribeAccessKeysWithToken(c, query)
//	case ACTION_ListAccessKeysByProject:
//		ListAccessKeysByProject(c, query)
//	case ACTION_CreateAccessKey:
//		CreateAccessKey(c, query)
//	case ACTION_DeleteAccessKey:
//		DeleteAccessKey(c, query)
//	default:
//		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"unsupport action",Data:""})
//		return
//	}
//	return
//}

func RegisterRequests(r *mux.Router) {
	r.HandleFunc("/api/v1/login", Login).Methods("POST")
	r.HandleFunc("/api/v1/yig/fetchsk", FetchSecretKey).Methods("POST")

	r.HandleFunc("/api/v1/account/create", CreateAccount).Methods("POST")
	r.HandleFunc("/api/v1/account/delete", DeleteAccount).Methods("POST")
	r.HandleFunc("/api/v1/account/describe", DescribeAccount).Methods("POST")
	r.HandleFunc("/api/v1/account/activate", ActivateAccount).Methods("POST")
	r.HandleFunc("/api/v1/account/deactivate", DeactivateAccount).Methods("POST")
	r.HandleFunc("/api/v1/account/list", ListAccounts).Methods("GET")

	r.HandleFunc("/api/v1/user/list", ListUsers).Methods("GET")
	r.HandleFunc("/api/v1/user/create", CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/user/delete", DeleteUser).Methods("POST")
	r.HandleFunc("/api/v1/user/describe", DescribeUser).Methods("POST")
	r.HandleFunc("/api/v1/user/listbyproject", ListUserByProject).Methods("POST")

	r.HandleFunc("/api/v1/project/list", ListProjects).Methods("GET")
	r.HandleFunc("/api/v1/project/create", CreateProject).Methods("POST")
	r.HandleFunc("/api/v1/project/delete", DeleteProject).Methods("POST")
	r.HandleFunc("/api/v1/project/describe", DescribeProject).Methods("POST")
	r.HandleFunc("/api/v1/project/listbyuser", ListProjectByUser).Methods("POST")
	r.HandleFunc("/api/v1/project/linkuser", LinkUserWithProject).Methods("POST")
	r.HandleFunc("/api/v1/project/unlinkuser", UnLinkUserWithProject).Methods("POST")
}

func ShowAPI(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			fmt.Println(err)
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		sm := strings.Join(m, ",")
		SRoute := map[string]string{"method": sm, "path": t}
		AllRoutes = append(AllRoutes, SRoute)
		return nil
	})

}

func TokenGc() {
	ticker := time.NewTicker(time.Minute*5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tokens, err := db.ListExpiredTokens()
			if err != nil {
				continue
			}
			for _, token := range tokens {
				db.RemoveToken(token.Token)
			}
		}
	}
}


func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
		AllowCredentials:   true,
	})
	helper.SetupConfig()
	r := mux.NewRouter()
	r.StrictSlash(true)
	n := negroni.New() // Includes some default middlewares
	logrusMiddleWare := negronilogrus.NewMiddleware()
	file, err := helper.OpenAccessLogFile()
	if err == nil {
		logrusMiddleWare.Logger.Out = file
		n.Use(logrusMiddleWare)
		defer file.Close()
	}
	if err := helper.CreatePidfile(helper.Config.PidFile); err != nil {
		fmt.Printf("can not create pid file %s\n", helper.Config.PidFile)
		return
	}
	defer helper.RemovePidfile(helper.Config.PidFile)
	RegisterRequests(r)
	ShowAPI(r)
	db.Db_Init()
	a := xormadapter.NewAdapter("mysql", helper.Config.RbacDataSource)
	helper.Enforcer = casbin.NewEnforcer("/etc/yig-iam/basic_model.conf", a)
	helper.Enforcer.LoadPolicy()
	roles := helper.Enforcer.GetAllRoles()
	if len(roles) == 0 {
	    helper.Logger.Println("roles number:", len(roles))
		helper.Casbin_init()
	}
	n.Use(c)
	n.Use(auth.Authorizer(helper.Enforcer))
	recovery := negroni.NewRecovery()
	recovery.Formatter = &negroni.HTMLPanicFormatter{}
	n.Use(recovery)
	n.UseHandler(r)
	go TokenGc()
	http.ListenAndServe(fmt.Sprintf(":%d", helper.Config.BindPort), n)
}

