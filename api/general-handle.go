package api
import (
	"gopkg.in/iris.v4"
	. "github.com/journeymidnight/yig-iam/api/datatype"
)
func EnvHandler(c *iris.Context) {
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:"{appName: manager}"})
}

func ApiHandle(c *iris.Context) {
	query := c.Get("queryRequest").(QueryRequest)

	switch query.Action {
	case ACTION_ConnectService:
		ConnectService(c, query)
	case ACTION_CreateAccount:
		CreateAccount(c, query)
	case ACTION_DeleteAccount:
		DeleteAccount(c, query)
	case ACTION_DescribeAccount:
		DescribeAccount(c, query)
	case ACTION_DeactivateAccount:
		DeactivateAccount(c, query)
	case ACTION_ActivateAccount:
		ActivateAccount(c, query)
	case ACTION_ListAccounts:
		ListAccounts(c, query)
	case ACTION_ListUsers:
		ListUsers(c, query)
	case ACTION_DescribeUser:
		DescribeUser(c, query)
	case ACTION_CreateUser:
		CreateUser(c, query)
	case ACTION_DeleteUser:
		DeleteUser(c, query)
	case ACTION_DescribeProject:
		DescribeProject(c, query)
	case ACTION_CreateProject:
		CreateProject(c, query)
	case ACTION_DeleteProject:
		DeleteProject(c, query)
	case ACTION_ListProjects:
		ListProjects(c, query)
	case ACTION_LinkUserWithProject:
		LinkUserWithProject(c, query)
	case ACTION_UnLinkUserWithProject:
		UnLinkUserWithProject(c, query)
	case ACTION_ListProjectByUser:
		ListProjectByUser(c, query)
	case ACTION_ListUserByProject:
		ListUserByProject(c, query)
	case ACTION_AddProjectService:
		AddProjectService(c, query)
	case ACTION_DelProjectService:
		DelProjectService(c, query)
	case ACTION_ListServiceByProject:
		ListServiceByProject(c, query)
	case ACTION_DescribeAccessKeys:
		DescribeAccessKeys(c, query)
	case ACTION_ListAccessKeysByProject:
		ListAccessKeysByProject(c, query)
	case ACTION_CreateAccessKey:
		CreateAccessKey(c, query)
	case ACTION_DeleteAccessKey:
		DeleteAccessKey(c, query)
	default:
		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"unsupport action",Data:""})
		return
	}
	return
}
