package store

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/journeymidnight/nier/src/helper"
	. "github.com/journeymidnight/nier/src/api/datatype"
)

func Casbin_init() {
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/login", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/create", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/delete", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/modify", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/user/roles", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/monitor/ceph", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/history/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/manage/snmp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/manage/snmp", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/manage/ntp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/manage/ntp", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/view/create", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/view/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/view/remove", "POST"})

	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/alerts", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/resolvedalerts", "GET"})

	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/catkeeper", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/catkeeper", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ROOT, "/api/v1/operate/addnewdisk", "POST"})


	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/user/login", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/monitor/ceph", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/history/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/manage/snmp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/manage/snmp", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/manage/ntp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/manage/ntp", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/view/create", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/view/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/view/remove", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/alerts", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/resolvedalerts", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/catkeeper", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/catkeeper", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_ADMIN, "/api/v1/operate/addnewdisk", "POST"})

	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/user/login", "POST"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/monitor/ceph", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/manage/snmp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/manage/ntp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/history/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/user/roles", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/view/list", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/manage/ntp", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/alerts", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/api/v1/resolvedalerts", "GET"})
	helper.Enforcer.AddPolicy([]string{ROLE_USER, "/catkeeper", "GET"})
	helper.Enforcer.AddRoleForUser("admin", ROLE_ROOT)
	helper.Enforcer.SavePolicy()
}

var Db backendstore

type backendstore interface {
	InsertUserRecord(userName string, password string, accountType string) error
	RemoveUserRecord(userName string) error
	ModifyUserRecord(userName string, password string, accountType string) error
	DescribeUserRecord(userName string) (UserRecord, error)
	ValidUserAndPassword(userName string, password string) (UserRecord, error)
	CheckUserExist(UserName string) (bool, error)
	InsertTokenRecord(Token string, UserName string, Type string) error
	RemoveTokenRecord(Token string) error
	ListExpiredTokens() ([]TokenRecord, error)
	GetTokenRecord(Token string) (TokenRecord, error)
	SearchExistedToken(userName string) (TokenRecord, error)
	ListUserRecords() ([]UserRecord, error)
	InsertSnmpRecord(snmpAddress string) error
	InsertNtpRecord(ntpAddress string) error
	GetSnmpRecord() (SnmpRecord, error)
	GetNtpRecord() (NtpRecord, error)
	GetHistoryRecord() ([]HistoryRecord, error)
	InsertOpRecord(userName, opType string) error
	GetAlertRecord() ([]AlertRecord, error)
	GetResolvedAlertRecord() ([]AlertRecord, error)
	Close()
}