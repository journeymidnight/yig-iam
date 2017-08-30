package helper

import (
	"github.com/hsluoyz/casbin"
	. "github.com/journeymidnight/yig-iam/api/datatype"
)

var Enforcer *casbin.Enforcer

func Casbin_init() {
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateProjectRole, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteProjectRole, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DescribeProjectRoles, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_GetAutogenkeysByProjectId, "ACCESS"})

	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateRegion, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteRegion, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ModifyRegionAttributes, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DescribeRegions, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateService, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteService, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DescribeServices, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ModifyServiceAttributes, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ModifyProjectAttributes, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ListProjects, "ACCESS"})

	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeactivateAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ActivateAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DescribeAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ListAccounts, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_ListUsers, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DescribeUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_CreateUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DeleteUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DescribeProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_GetLinkedProjectsByAccount, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_CreateProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ModifyProjectAttributes, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_DeleteProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, API_ListProjects, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_LinkUserWithProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_UnLinkUserWithProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_ListProjectByUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_ListUserByProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_AddProjectService, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DelProjectService, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_ListServiceByProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_ListAccessKeysByProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_CreateAccessKey, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DeleteAccessKey, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, API_DescribeAccessKeysWithToken, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_USER, API_DescribeUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_USER, API_DescribeProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_USER, API_ListAccessKeysByProject, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_USER, API_ListProjectByUser, "ACCESS"})
	Enforcer.AddPolicy([]string{ROLE_USER, API_ListServiceByProject, "ACCESS"})

	Enforcer.AddRoleForUser("root", ROLE_ROOT)
	Enforcer.SavePolicy()
}
