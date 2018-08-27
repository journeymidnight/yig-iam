package helper

import (
	. "github.com/journeymidnight/yig-iam/api/datatype"
)

func Casbin_init () {
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/create", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/delete", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/describe", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/activate", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/deactivate", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/account/list", "GET"})
	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/update", "POST"})

	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/list", "GET"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/create", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/delete", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/describe", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/listbyproject", "POST"}) //both owned by account and user
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/user/update", "POST"})

	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/list", "GET"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/create", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/delete", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/describe", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/listbyuser", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/linkuser", "POST"})
	Enforcer.AddPolicy([]string{ROLE_ACCOUNT, "/api/v1/project/unlinkuser", "POST"})

	Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/user/update", "POST"})
	Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/project/list", "GET"})

	Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/project/list", "GET"})

	Enforcer.AddRoleForUser("root", ROLE_ROOT)
	Enforcer.AddRoleForUser("admin", ROLE_ACCOUNT)
	Enforcer.SavePolicy()
}
