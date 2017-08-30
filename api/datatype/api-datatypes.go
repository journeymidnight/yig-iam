package datatype

const (
	ROLE_ROOT    = "ROOT"
	ROLE_ACCOUNT = "ACCOUNT"
	ROLE_USER    = "USER"
)

const (
	ACT_ACCESS = "ACCESS"
)

const (
	API_CreateProjectRole          = "CreateProjectRole"
	API_DeleteProjectRole          = "DeleteProjectRole"
	API_DescribeProjectRoles       = "DescribeProjectRoles"
	API_GetLinkedProjectsByAccount = "GetLinkedProjectsByAccount"
	API_GetAutogenkeysByProjectId  = "GetAutogenkeysByProjectId"

	//supermanager apis
	API_CreateRegion            = "CreateRegion"
	API_DeleteRegion            = "DeleteRegion"
	API_ModifyRegionAttributes  = "ModifyRegionAttributes"
	API_DescribeRegions         = "DescribeRegions"
	API_CreateService           = "CreateService"
	API_DeleteService           = "DeleteService"
	API_ModifyServiceAttributes = "ModifyServiceAttributes"
	API_DescribeServices        = "DescribeServices"

	API_CreateAccount               = "CreateAccount"
	API_DeleteAccount               = "DeleteAccount"
	API_DescribeAccount             = "DescribeAccount"
	API_DeactivateAccount           = "DeactivateAccount"
	API_ActivateAccount             = "ActivateAccount"
	API_ListAccounts                = "ListAccounts"
	API_ListUsers                   = "ListUsers"
	API_DescribeUser                = "DescribeUser"
	API_CreateUser                  = "CreateUser"
	API_DeleteUser                  = "DeleteUser"
	API_DescribeProject             = "DescribeProject"
	API_CreateProject               = "CreateProject"
	API_ModifyProjectAttributes     = "ModifyProjectAttributes"
	API_DeleteProject               = "DeleteProject"
	API_ListProjects                = "ListProjects"
	API_LinkUserWithProject         = "LinkUserWithProject"
	API_UnLinkUserWithProject       = "UnLinkUserWithProject"
	API_ListProjectByUser           = "ListProjectByUser"
	API_ListUserByProject           = "ListUserByProject"
	API_AddProjectService           = "AddProjectService"
	API_DelProjectService           = "DelProjectService"
	API_ListServiceByProject        = "ListServiceByProject"
	API_DescribeAccessKeys          = "DescribeAccessKeys" //priviate api for internal system such as yig
	API_ListAccessKeysByProject     = "ListAccessKeysByProject"
	API_DescribeAccessKeysWithToken = "DescribeAccessKeysWithToken" //priviate api for internal system such as yig
	API_CreateAccessKey             = "CreateAccessKey"
	API_DeleteAccessKey             = "DeleteAccessKey"
)

const (
	ACTION_CreateProjectRole           = "CreateProjectRole"
	ACTION_DeleteProjectRole           = "DeleteProjectRole"
	ACTION_DescribeProjectRoles        = "DescribeProjectRoles"
	ACTION_GetLinkedProjectsByAccount  = "GetLinkedProjectsByAccount"
	ACTION_GetAutogenkeysByProjectId   = "GetAutogenkeysByProjectId"
	ACTION_CreateRegion                = "CreateRegion"
	ACTION_DeleteRegion                = "DeleteRegion"
	ACTION_DescribeRegions             = "DescribeRegions"
	ACTION_ModifyRegionAttributes      = "ModifyRegionAttributes"
	ACTION_CreateService               = "CreateService"
	ACTION_DeleteService               = "DeleteService"
	ACTION_ModifyServiceAttributes     = "ModifyServiceAttributes"
	ACTION_DescribeServices            = "DescribeServices"
	ACTION_ConnectService              = "ConnectService"
	ACTION_CreateAccount               = "CreateAccount"
	ACTION_DeleteAccount               = "DeleteAccount"
	ACTION_DescribeAccount             = "DescribeAccount"
	ACTION_DeactivateAccount           = "DeactivateAccount"
	ACTION_ActivateAccount             = "ActivateAccount"
	ACTION_ListAccounts                = "ListAccounts"
	ACTION_ListUsers                   = "ListUsers"
	ACTION_DescribeUser                = "DescribeUser"
	ACTION_CreateUser                  = "CreateUser"
	ACTION_DeleteUser                  = "DeleteUser"
	ACTION_DescribeProject             = "DescribeProject"
	ACTION_CreateProject               = "CreateProject"
	ACTION_ModifyProjectAttributes     = "ModifyProjectAttributes"
	ACTION_DeleteProject               = "DeleteProject"
	ACTION_ListProjects                = "ListProjects"
	ACTION_LinkUserWithProject         = "LinkUserWithProject"
	ACTION_UnLinkUserWithProject       = "UnLinkUserWithProject"
	ACTION_ListProjectByUser           = "ListProjectByUser"
	ACTION_ListUserByProject           = "ListUserByProject"
	ACTION_AddProjectService           = "AddProjectService"
	ACTION_DelProjectService           = "DelProjectService"
	ACTION_ListServiceByProject        = "ListServiceByProject"
	ACTION_DescribeAccessKeys          = "DescribeAccessKeys"          //priviate api for internal system such as yig
	ACTION_DescribeAccessKeysWithToken = "DescribeAccessKeysWithToken" //priviate api for internal system such as yig
	ACTION_ListAccessKeysByProject     = "ListAccessKeysByProject"
	ACTION_CreateAccessKey             = "CreateAccessKey"
	ACTION_DeleteAccessKey             = "DeleteAccessKey"
	//losapis
	ACTION_LOS_GETS3DOMAIN    = "GetS3Domain"
	ACTION_LOS_LISTBUCKETS    = "ListBuckets"
	ACTION_LOS_DELETEBUCKET   = "DeleteBucket"
	ACTION_LOS_GETBUCKETSTATS = "GetBucketStats"
	ACTION_LOS_CREATEBUCKET   = "CreateBucket"
	ACTION_LOS_PUTCORS        = "PutCors"
)

type QueryRequest struct {
	Action      string   `json:"action"`
	AccountId   string   `json:"accountId,omitempty"`
	UserName    string   `json:"user,omitempty"`
	KeyName     string   `json:"keyName,omitempty"`
	ProjectId   string   `json:"projectId,omitempty"`
	ProjectName string   `json:"projectName,omitempty"`
	ProjectIds  []string `json:"projects,omitempty"`
	Password    string   `json:"password,omitempty"`
	Description string   `json:"description,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Type        string   `json:"type,omitempty"`
	Email       string   `json:"email,omitempty"`
	Service     string   `json:"Service,omitempty"`
	Quota       string   `json:"quota,omitempty"`
	Token       string   `json:"token,omitempty"`
	AccessKey   string   `json:"accessKey,omitempty"`
	AccessKeys  []string `json:"accessKeys,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
	Bucket      string   `json:"bucket,omitempty"`

	// region related
	RegionName string `json:"regionName,omitempty"`
	RegionId   string `json:"regionId,omitempty"`

	// service related
	ServiceId string `json:"serviceId,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	PublicUrl string `json:"publicurl,omitempty"`

	// role releated
	Role int `json:"role,omitempty"`
}

type QueryRespAll struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	RetCode int         `json:"retCode"`
}

/************compatible for YIG**************/

type AccessKeyItem struct {
	ProjectId    string `json:"projectId"`
	Name         string `json:"name"`
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
	Status       string `json:"status"`
	Updated      string `json:"updated"`
	Created      string `json:"created"`
	Description  string `json:"description"`
}

type DescribeKeysResp struct {
	Limit        int             `json:"limit"`
	Total        int             `json:"total"`
	Offset       int             `json:"offset"`
	AccessKeySet []AccessKeyItem `json:"accessKeySet"`
}

/***********************************/

type QueryResponse struct {
	RetCode int         `json:"retCode"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ConnectServiceResponse struct {
	Token      string                `json:"token"`
	Type       string                `json:"type"`
	AccountId  string                `json:"accountId"`
	ProjectSet []LinkedProjectRecord `json:"projectSet"`
}

type UserRecord struct {
	UserName    string
	Password    string
	Type        string
	Email       string
	DisplayName string
	AccountId   string
	Status      string
	Created     string
	Updated     string
}

type ProjectRecord struct {
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	AccountId   string `json:"accountId"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

type ListProjectResp struct {
	Limit    int             `json:"limit"`
	Total    int             `json:"total"`
	Offset   int             `json:"offset"`
	Projects []ProjectRecord `json:"projectSet"`
}

type UserProjectRecord struct {
	UserName  string `json:"userName"`
	ProjectId string `json:"projectId"`
	Created   string `json:"created"`
}

type ProjectServiceRecord struct {
	ProjectId string `json:"projectId"`
	Service   string `json:"service"`
	AccountId string `json:"accountId"`
	Created   string `json:"created"`
}

type AkSkRecord struct {
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
	ProjectId    string `json:"projectId"`
	AccountId    string `json:"accountId"`
	KeyName      string `json:"keyName"`
	Created      string `json:"created"`
	Description  string `json:"description"`
}

type TokenRecord struct {
	Token     string `json:"token"`
	UserName  string `json:"userName"`
	AccountId string `json:"accountId"`
	Type      string `json:"type"`
	Created   string `json:"created"`
	Expired   string `json:"expired"`
}

type RegionRecord struct {
	RegionId   string `json:"regionId"`
	RegionName string `json:"regionName"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
}

type ListRegionResp struct {
	Limit   int            `json:"limit"`
	Total   int            `json:"total"`
	Offset  int            `json:"offset"`
	Regions []RegionRecord `json:"regionSet"`
}

type ServiceRecord struct {
	ServiceId string `json:"serviceId"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Endpoint  string `json:"endpoint"`
	RegionId  string `json:"regionId"`
}

type ListServiceResp struct {
	Limit    int             `json:"limit"`
	Total    int             `json:"total"`
	Offset   int             `json:"offset"`
	Services []ServiceRecord `json:"serviceSet"`
}

type ProjectRoleRecord struct {
	UserId      string `json:"userId"`
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Role        int    `json:"role"`
	Created     string `json:"created"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

type ListProjectRoleResp struct {
	Limit   int                 `json:"limit"`
	Total   int                 `json:"total"`
	Offset  int                 `json:"offset"`
	RoleSet []ProjectRoleRecord `json:"roleSet"`
}

type LinkedProjectRecord struct {
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
}

type LinkedProjectsResp struct {
	ProjectSet []LinkedProjectRecord `json:"projectSet"`
}
