package datatype

import "time"

const (
	ROLE_ROOT = "ROOT"
	ROLE_ACCOUNT = "ACCOUNT"
	ROLE_USER = "USER"
	REQUEST_TOKEN_KEY = "TOKEN"
	USER_STATUS_INACTIVE	= "inactive"
	USER_STATUS_ACTIVE	= "active"
	PROJECT_STATUS_INACTIVE	= "inactive"
	PROJECT_STATUS_ACTIVE	= "active"
	PROJECT_STATUS_DELETED	= "deleted"
	KEY_STATUS_ENABLE = "enable"
	KEY_STATUS_DISABLE = "disable"
	ACL_RW = "RW"
	ACL_RO = "RDONLY"
	PUBLIC_PROJECT = "public"
	PRIVATE_PROJECT = "private"
)

const (
	ACT_ACCESS = "ACCESS"
)

const (
	ACTION_ConnectService = "ConnectService"
	ACTION_CreateAccount = "CreateAccount"
	ACTION_DeleteAccount = "DeleteAccount"
	ACTION_DescribeAccount = "DescribeAccount"
	ACTION_DeactivateAccount = "DeactivateAccount"
	ACTION_ActivateAccount = "ActivateAccount"
	ACTION_ListAccounts = "ListAccounts"
	ACTION_ListUsers = "ListUsers"
	ACTION_DescribeUser = "DescribeUser"
	ACTION_CreateUser = "CreateUser"
	ACTION_DeleteUser = "DeleteUser"
	ACTION_DescribeProject = "DescribeProject"
	ACTION_CreateProject = "CreateProject"
	ACTION_DeleteProject = "DeleteProject"
	ACTION_ListProjects = "ListProjects"
	ACTION_LinkUserWithProject = "LinkUserWithProject"
	ACTION_UnLinkUserWithProject = "UnLinkUserWithProject"
	ACTION_ListProjectByUser = "ListProjectByUser"
	ACTION_ListUserByProject = "ListUserByProject"
	ACTION_AddProjectService = "AddProjectService"
	ACTION_DelProjectService = "DelProjectService"
	ACTION_ListServiceByProject = "ListServiceByProject"
	ACTION_DescribeAccessKeys = "DescribeAccessKeys" //priviate api for internal system such as yig
	ACTION_DescribeAccessKeysWithToken = "DescribeAccessKeysWithToken" //priviate api for internal system such as yig
	ACTION_ListAccessKeysByProject = "ListAccessKeysByProject"
	ACTION_CreateAccessKey = "CreateAccessKey"
	ACTION_DeleteAccessKey = "DeleteAccessKey"
)

type JsonTime time.Time
func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"`+time.Time(j).Format("2006-01-02 15:04:05")+`"`), nil
}

type Token struct {
	Token     string     `json:"token" xorm:"'token' varchar(40) pk"`
	UserName  string     `json:"userName"  xorm:"'userName' varchar(50) notnull"`
	UserId  string     `json:"userId"  xorm:" 'userId' pk varchar(20) notnull"`
	AccountId  string     `json:"accountId"  xorm:" 'accountId' varchar(20) notnull"`
	Type      string     `json:"type"  xorm:"'type' varchar(10)  notnull"`
	CreatedAt  JsonTime     `json:"created"  xorm:" 'created' created"`
}


type User struct {
	UserName 	string    `json:"userName" xorm:" 'userName' varchar(50) pk notnull"`
	Password    string     `json:"password" xorm:" 'password' varchar(20) notnull"`
	UserId    string     `json:"userId" xorm:" 'userId' varchar(20) notnull"`
	Type        string     `json:"type"  xorm:" 'type' varchar(10) notnull"`
	Email        string     `json:"email"  xorm:" 'email' varchar(50) DEFAULT NULL"`
	DisplayName  string     `json:"displayName"  xorm:" 'displayName' varchar(50) DEFAULT NULL"`
	AccountId  string     `json:"accountId"  xorm:" 'accountId' varchar(20) notnull"`
	Status  string     `json:"status"  xorm:" 'status' varchar(10) notnull"`
	CreatedAt  JsonTime     `json:"created"  xorm:" 'created' created"`
	UpdatedAt  JsonTime     `json:"updated"  xorm:" 'updated' updated"`
}

type Project struct {
	ProjectId 	string    `json:"projectId" xorm:" 'projectId' pk varchar(20) notnull"`
	ProjectName    string     `json:"projectName" xorm:" 'projectName' varchar(50) notnull"`
	ProjectType    string     `json:"projectType" xorm:" 'projectType' varchar(20) notnull"`
	AccountId  string     `json:"accountId"  xorm:" 'accountId' varchar(20) notnull"`
	OwnerId  string     `json:"ownerId"  xorm:" 'ownerId' varchar(20) notnull"`
	Description  string     `json:"description"  xorm:" 'description' varchar(50) DEFAULT NULL"`
	Status  string     `json:"status"  xorm:" 'status' varchar(10) notnull"`
	CreatedAt  JsonTime     `json:"created"  xorm:" 'created' created"`
	UpdatedAt  JsonTime     `json:"updated"  xorm:" 'updated' updated"`
}

type UserProject struct {
	AccessKey 	string    `json:"accessKey" xorm:" 'accessKey' varchar(20) notnull"`
	AccessSecret    string     `json:"accessSecret" xorm:" 'accessSecret' varchar(40) notnull"`
	Status  string     `json:"status"  xorm:" 'status' varchar(10) notnull"`
	Acl    string     `json:"acl" xorm:" 'acl' varchar(20) notnull"`
	ProjectId  string     `json:"projectId"  xorm:" 'projectId' pk varchar(20) notnull"`
	UserId  string     `json:"userId"  xorm:" 'userId' pk varchar(20) notnull"`
	AccountId  string     `json:"accountId"  xorm:" 'accountId' varchar(20) notnull"`
	CreatedAt  JsonTime     `json:"created"  xorm:" 'created' created"`
	UpdatedAt  JsonTime     `json:"updated"  xorm:" 'updated' updated"`
}

type QueryRequest struct{
	Acl string `json:"acl,omitempty"`
	UserName string `json:"userName,omitempty"`
	UserId string `json:"userId,omitempty"`
	KeyName string `json:"keyName,omitempty"`
	ProjectId string `json:"projectId,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
	ProjectIds []string `json:"projects,omitempty"`
	Password string `json:"password,omitempty"`
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Type string `json:"type,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
	AccessKey string `json:"accessKey,omitempty"`
	AccessKeys []string `json:"accessKeys,omitempty"`
	Limit int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type FetchAccessKeysResp struct {
	AccessKeySet []AccessKeyItem `json:"accessKeySet"`
}

/************compatible for YIG**************/

type AccessKeyItem struct {
	ProjectId    string `json:"projectId"`
	Name         string `json:"name"`
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
	Acl          string `json:"acl"`
	Status       string `json:"status"`
	Updated      string `json:"updated"`
}

type QueryResp struct {
	Limit        int             `json:"limit"`
	Total        int             `json:"total"`
	Offset       int             `json:"offset"`
	AccessKeySet []AccessKeyItem `json:"accessKeySet"`
}

type QueryRespAll struct {
	Message string    `json:"message"`
	Data    interface{} `json:"data"`
	RetCode int       `json:"retCode"`
}

/***********************************/

type LoginResponse struct {
	Token string `json:"token"`
	Type string `json:"type"`
	UserId string `json:"userId"`
	AccountId string `json:"accountId"`
}

//type UserRecord struct {
//	UserName string
//	Password string
//	Type     string
//	Email    string
//	DisplayName string
//	AccountId string
//	Status   string
//	Created  string
//	Updated  string
//}
//
//type ProjectRecord struct {
//	ProjectId string `json:"projectId"`
//	ProjectName string `json:"projectName"`
//	ProjectType string `json:"projectType"`
//	AccountId string `json:"accountId"`
//	Description string `json:"description"`
//	Status string `json:"status"`
//	Created  string `json:"created"`
//	Updated  string `json:"updated"`
//}

type ListProjectResp struct {
	UserProject `xorm:"extends"`
	ProjectName    string `json:"projectName" xorm:" 'projectName' varchar(50) notnull"`
	ProjectType    string `json:"projectType" xorm:" 'projectType' varchar(20) notnull"`
}

func (ListProjectResp) TableName() string {
	return "user_project"
}

type ListUserResp struct {
	UserProject `xorm:"extends"`
	UserName    string `json:"userName" xorm:" 'userName' varchar(50) pk notnull"`
	DisplayName    string `json:"displayName"  xorm:" 'displayName' varchar(50) DEFAULT NULL"`
}

func (ListUserResp) TableName() string {
	return "user_project"
}

//type UserProjectRecord struct {
//	UserName string `json:"userName"`
//	ProjectId string `json:"projectId"`
//	Created  string `json:"created"`
//}
//
////type ProjectServiceRecord struct {
////	ProjectId string `json:"projectId"`
////	Service string `json:"service"`
////	AccountId string `json:"accountId"`
////	Created  string `json:"created"`
////}
//
//type AkSkRecord struct {
//	AccessKey string `json:"accessKey"`
//	AccessSecret string `json:"accessSecret"`
//	ProjectId string `json:"projectId"`
//	AccountId string `json:"accountId"`
//	KeyName string `json:"keyName"`
//	Created  string `json:"created"`
//	Description string `json:"description"`
//}
//
//type TokenRecord struct {
//	Token string `json:"token"`
//	UserName string `json:"userName"`
//	AccountId string `json:"accountId"`
//	Type  string `json:"type"`
//	Created  string `json:"created"`
//	Expired  string `json:"expired"`
//}
