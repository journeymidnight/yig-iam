package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/helper"
)

var Db *sql.DB

const (
	User_TYPE_ADMIN   = 0
	User_TYPE_ACCOUNT = 1
	User_TYPE_USER    = 2
)

const TimeFormat = "2006-01-02T15:04:05Z07:00"

func CreateDbConnection() *sql.DB {
	conn, err := sql.Open("mysql", helper.CONFIG.DatabaseConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}
	helper.Logger.Println(5, "Connected to database")
	return conn
}

func checkDbTables() {
	return
}

// projectRole table
func InsertProjectRoleRecord(projectId, projectName, userid string, role int) error {
	created := time.Now().Format(TimeFormat)
	_, err := Db.Exec("insert into ProjectUser values( null, ?, ?, ?, ?, ?)", userid, projectId, projectName, role, created)
	if err != nil {
		helper.Logger.Println(5, "Error add ProjectUser", userid, projectId, projectName, role, created, err.Error())
	}
	return err
}

func GetLinkedProjects(account string) ([]LinkedProjectRecord, error) {
	var records []LinkedProjectRecord
	rows, err := Db.Query("select * from ProjectUser where user_id=(?)", account)
	if err != nil {
		helper.Logger.Println(5, "Error GetLinkedProjects: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record LinkedProjectRecord
		var index int
		var userid, role, created string
		if err := rows.Scan(&index,
			&userid,
			&record.ProjectId,
			&record.ProjectName,
			&role,
			&created); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func ListProjectRoleRecordsByProjectId(projectid string) ([]ProjectRoleRecord, error) {
	var records []ProjectRoleRecord
	rows, err := Db.Query("select * from ProjectUser where project_id=(?)", projectid)
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record ProjectRoleRecord
		var someid int
		if err := rows.Scan(
			&someid,
			&record.UserId,
			&record.ProjectId,
			&record.ProjectName,
			&record.Role,
			&record.Created); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func RemoveProjectRoleRecord(ProjectId string, AccountId string) error {
	_, err := Db.Exec("delete from ProjectUser where project_id=(?) and user_id=(?)", ProjectId, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error remove projectuser", ProjectId, err.Error())
	}
	return err
}
func InsertRegionRecord(RegionId string, RegionName string) error {
	status := "active"
	created := time.Now().Format(TimeFormat)
	updated := created
	_, err := Db.Exec("insert into Region values( ?, ?, ?, ?, ?)", RegionId, RegionName, status, created, updated)
	if err != nil {
		helper.Logger.Println(5, "Error add region", RegionId, RegionName, status, created, updated, err.Error())
	}
	return err
}

func ListRegionRecords() ([]RegionRecord, error) {
	var records []RegionRecord
	rows, err := Db.Query("select * from Region")
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record RegionRecord
		if err := rows.Scan(
			&record.RegionId,
			&record.RegionName,
			&record.Status,
			&record.Created,
			&record.Updated); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func UpdateRegionRecord(id, name string) error {
	updated := time.Now().Format(TimeFormat)
	_, err := Db.Exec("update Region set regionName=(?), updated=(?) where regionId=(?)", name, id, updated)
	if err != nil {
		helper.Logger.Println(5, "Error modify project", id, err.Error())
	}
	return err
}

func RemoveRegionRecord(id string) error {
	_, err := Db.Exec("delete from Region where regionId=(?)", id)
	if err != nil {
		helper.Logger.Println(5, "Error remove region", id, err.Error())
	}

	return err
}

func InsertServiceRecord(serviceid, url, endpoint, regionId string) error {
	created := time.Now().Format(TimeFormat)
	updated := created
	_, err := Db.Exec("insert into Service values( ?, ?, ?, ?, ?)", serviceid, created, updated, regionId, endpoint)
	if err != nil {
		helper.Logger.Println(5, "Error add service", serviceid, created, updated, regionId, endpoint, err.Error())
	}
	return err
}

func ListSerivceRecords() ([]ServiceRecord, error) {
	var records []ServiceRecord
	rows, err := Db.Query("select * from Service")
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record ServiceRecord
		if err := rows.Scan(
			&record.ServiceId,
			&record.Created,
			&record.Updated,
			&record.RegionId,
			&record.Endpoint); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func UpdateServiceRecord(id, endpoint string) error {
	updated := time.Now().Format(TimeFormat)
	_, err := Db.Exec("update Service set Endpoint=(?), updated=(?) where serviceId=(?)", endpoint, updated, id)
	if err != nil {
		helper.Logger.Println(5, "Error modify service", id, err.Error())
	}
	return err
}

func RemoveServiceRecord(id string) error {
	_, err := Db.Exec("delete from Service where serviceId=(?)", id)
	if err != nil {
		helper.Logger.Println(5, "Error remove service", id, err.Error())
	}

	return err
}

func RemoveAccountId(AccountId string) error {
	_, err := Db.Exec("delete from User where accountId=(?) and type='ACCOUNT'", AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error remove account", AccountId, err.Error())
	}
	helper.Logger.Println(5, "DeleteAccount:", err)
	return err
}

func DescribeAccount(AccountId string) (UserRecord, error) {
	var record UserRecord
	err := Db.QueryRow("select * from User where accountId=(?) and type='ACCOUNT'", AccountId).Scan(
		&record.UserName,
		&record.Password,
		&record.Type,
		&record.Email,
		&record.DisplayName,
		&record.AccountId,
		&record.Status,
		&record.Created,
		&record.Updated)
	record.Password = ""
	return record, err
}

func DeactivateAccount(AccountId string) error {
	_, err := Db.Exec("update User set status='inactive' where accountId=(?) and type='ACCOUNT'", AccountId)
	return err
}

func ActivateAccount(AccountId string) error {
	_, err := Db.Exec("update User set status='active' where accountId=(?) and type='ACCOUNT'", AccountId)
	return err
}

func ListAccountRecords() ([]UserRecord, error) {
	var records []UserRecord
	rows, err := Db.Query("select * from User where type=(?)", ROLE_ACCOUNT)
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record UserRecord
		if err := rows.Scan(
			&record.UserName,
			&record.Password,
			&record.Type,
			&record.Email,
			&record.DisplayName,
			&record.AccountId,
			&record.Status,
			&record.Created,
			&record.Updated); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		record.Password = ""
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func InsertUserRecord(userName string, password string, accountType string,
	email string, displayName string, accountId string) error {
	status := "active"
	created := time.Now().Format(TimeFormat)
	updated := created
	_, err := Db.Exec("insert into User values( ?, ?, ?, ?, ?, ?, ?, ?, ? )", userName, password, accountType,
		email, displayName, accountId, status, created, updated)
	if err != nil {
		helper.Logger.Println(5, "Error add account", userName, password, accountType,
			email, displayName, accountId, status, created, updated, err.Error())
	}
	return err
}

func RemoveUserRecord(userName string, accountId string) error {
	_, err := Db.Exec("delete from User where userName=(?) and accountId=(?)", userName, accountId)
	if err != nil {
		helper.Logger.Println(5, "Error remove user", userName, err.Error())
	}
	return err
}

func DescribeUserRecord(userName string, accountId string) (UserRecord, error) {
	var record UserRecord
	err := Db.QueryRow("select * from User where userName=(?) and accountId=(?)", userName, accountId).Scan(&record.UserName,
		&record.Password,
		&record.Type,
		&record.Email,
		&record.DisplayName,
		&record.AccountId,
		&record.Status,
		&record.Created,
		&record.Updated)
	record.Password = ""
	return record, err
}

func ValidEmailAndPassword(email, password string) (UserRecord, error) {
	var record UserRecord
	err := Db.QueryRow("select * from User where email=(?) and password=(?)", email, password).Scan(
		&record.UserName,
		&record.Password,
		&record.Type,
		&record.Email,
		&record.DisplayName,
		&record.AccountId,
		&record.Status,
		&record.Created,
		&record.Updated)
	record.Password = ""
	return record, err
}

func ValidUserAndPassword(userName string, password string) (UserRecord, error) {
	var record UserRecord
	err := Db.QueryRow("select * from User where userName=(?) and password=(?)", userName, password).Scan(
		&record.UserName,
		&record.Password,
		&record.Type,
		&record.Email,
		&record.DisplayName,
		&record.AccountId,
		&record.Status,
		&record.Created,
		&record.Updated)
	record.Password = ""
	return record, err
}

func ListUserRecords(accountId string) ([]UserRecord, error) {
	var records []UserRecord
	rows, err := Db.Query("select * from User where accountId=(?) and type=(?)", accountId, ROLE_USER)
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record UserRecord
		if err := rows.Scan(
			&record.UserName,
			&record.Password,
			&record.Type,
			&record.Email,
			&record.DisplayName,
			&record.AccountId,
			&record.Status,
			&record.Created,
			&record.Updated); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		record.Password = ""
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func InsertProjectRecord(ProjectId string, ProjectName string, AccountId string, Description string) error {
	status := "active"
	created := time.Now().Format(TimeFormat)
	updated := created
	_, err := Db.Exec("insert into Project values( ?, ?, ?, ?, ?, ?, ? )", ProjectId, ProjectName, AccountId, Description, status, created, updated)
	if err != nil {
		helper.Logger.Println(5, "Error add project", ProjectId, ProjectName, AccountId, Description, status, created, updated, err.Error())
	}
	return err
}

func RemoveProjectRecord(ProjectId string, AccountId string) error {
	_, err := Db.Exec("delete from Project where projectId=(?) and accountId=(?)", ProjectId, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error remove user", ProjectId, err.Error())
	}
	return err
}

func UpdateProjectRecord(ProjectId, ProjectName, Description string) error {
	updated := time.Now().Format(TimeFormat)
	_, err := Db.Exec("update Project set projectName=(?), description=(?), updated=(?) where projectId=(?)", ProjectName, Description, updated, ProjectId)
	if err != nil {
		helper.Logger.Println(5, "Error modify project", ProjectId, err.Error())
	}
	return err
}

func DescribeProjectRecord(ProjectId string, AccountId string) (ProjectRecord, error) {
	var record ProjectRecord
	err := Db.QueryRow("select * from Project where projectId=(?) and accountId=(?)", ProjectId, AccountId).Scan(&record.ProjectId,
		&record.ProjectName,
		&record.AccountId,
		&record.Description,
		&record.Status,
		&record.Created,
		&record.Updated)
	return record, err
}

func DescribeProjectRecordByProjectId(ProjectId string) (ProjectRecord, error) {
	var record ProjectRecord
	err := Db.QueryRow("select * from Project where projectId=(?)", ProjectId).Scan(&record.ProjectId,
		&record.ProjectName,
		&record.AccountId,
		&record.Description,
		&record.Status,
		&record.Created,
		&record.Updated)
	return record, err
}

func ListProjectRecords(accountId string) ([]ProjectRecord, error) {
	var records []ProjectRecord
	rows, err := Db.Query("select * from Project where accountId=(?)", accountId)
	if err != nil {
		helper.Logger.Println(5, "Error querying idle executors: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record ProjectRecord
		if err := rows.Scan(
			&record.ProjectId,
			&record.ProjectName,
			&record.AccountId,
			&record.Description,
			&record.Status,
			&record.Created,
			&record.Updated); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func InsertUserProjectRecord(ProjectId string, UserName string) error {
	created := time.Now().Format(TimeFormat)
	_, err := Db.Exec("insert into UserProject values( null, ?, ?, ? )", ProjectId, UserName, created)
	if err != nil {
		helper.Logger.Println(5, "Error add project", ProjectId, UserName, created, err.Error())
	}
	return err
}

func RemoveUserProjectRecord(ProjectId string, UserName string) error {
	_, err := Db.Exec("delete from UserProject where projectId=(?) and userName=(?)", ProjectId, UserName)
	if err != nil {
		helper.Logger.Println(5, "Error remove user", ProjectId, err.Error())
	}
	return err
}

func ListUserProjectRecordByUser(UserName string) ([]UserProjectRecord, error) {
	var records []UserProjectRecord
	rows, err := Db.Query("select * from UserProject where userName=(?)", UserName)
	if err != nil {
		helper.Logger.Println(5, "Error ListUserProjectRecordByUser: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record UserProjectRecord
		var index int
		if err := rows.Scan(&index, &record.ProjectId, &record.UserName, &record.Created); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func ListUserProjectRecordByProject(ProjectId string) ([]UserProjectRecord, error) {
	var records []UserProjectRecord
	rows, err := Db.Query("select * from UserProject where projectId=(?)", ProjectId)
	if err != nil {
		helper.Logger.Println(5, "Error ListUserProjectRecordByUser: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record UserProjectRecord
		var index int
		if err := rows.Scan(&index, &record.ProjectId, &record.UserName, &record.Created); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func InsertProjectServiceRecord(ProjectId string, Service string, AccountId string) error {
	created := time.Now().Format(TimeFormat)
	_, err := Db.Exec("insert into ProjectService values( null, ?, ?, ?, ? )", ProjectId, Service, AccountId, created)
	if err != nil {
		helper.Logger.Println(5, "Error InsertProjectServiceRecord", ProjectId, Service, AccountId, created, err.Error())
	}
	return err
}

func RemoveProjectServiceRecord(ProjectId string, Service string, AccountId string) error {
	_, err := Db.Exec("delete from ProjectService where projectId=(?) and service=(?) and accountId=(?)", ProjectId, Service, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error RemoveProjectServiceRecord", ProjectId, Service, AccountId, err.Error())
	}
	return err
}

func ListProjectServiceRecordByProject(ProjectId string, AccountId string) ([]ProjectServiceRecord, error) {
	var records []ProjectServiceRecord
	rows, err := Db.Query("select * from ProjectService where projectId=(?) and accountId=(?)", ProjectId, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error ListProjectServiceRecordByProject: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record ProjectServiceRecord
		var index int
		if err := rows.Scan(&index, &record.ProjectId, &record.Service, &record.AccountId, &record.Created); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return records, err
}

func CheckAccountIdExist(AccountId string) (bool, error) {
	var count int
	err := Db.QueryRow("select count(*) from User where AccountId=(?)", AccountId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckUserExist(UserName string) (bool, error) {
	var count int
	err := Db.QueryRow("select count(*) from User where userName=(?)", UserName).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}

}

func InsertTokenRecord(Token string, UserName string, AccountId string, Type string) error {
	created := time.Now().Format(TimeFormat)
	expired := time.Now().Add(time.Duration(helper.CONFIG.TokenExpire * 1000000000)).Format(TimeFormat)
	_, err := Db.Exec("insert into Token values( ?, ?, ?, ?, ?, ? )", Token, UserName, AccountId, Type, created, expired)
	if err != nil {
		helper.Logger.Println(5, "Error InsertTokenRecord", Token, UserName, AccountId, Type, created, expired, err.Error())
	}
	return err
}

func GetTokenRecord(Token string) (TokenRecord, error) {
	var record TokenRecord
	err := Db.QueryRow("select * from Token where token=(?)", Token).Scan(&record.Token,
		&record.UserName,
		&record.AccountId,
		&record.Type,
		&record.Created,
		&record.Expired)
	return record, err
}

func InsertAkSkRecord(AccessKey string, SecretKey string, ProjectId string, AccountId string, KeyName string, Description string, isAutoGenerated bool) error {
	created := time.Now().Format(TimeFormat)
	_, err := Db.Exec("insert into AkSk values( ?, ?, ?, ?, ?, ?, ?, ?)", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created, Description, isAutoGenerated)
	if err != nil {
		helper.Logger.Println(5, "Error InsertAkSkRecord ", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created, Description, isAutoGenerated, err.Error())
	}
	return err
}

func IfAKExisted(AccessKey string) bool {
	var record AkSkRecord
	var isautogen bool
	err := Db.QueryRow("select * from User where accessKey=(?)", AccessKey).Scan(
		&record.AccessKey,
		&record.AccessSecret,
		&record.ProjectId,
		&record.AccountId,
		&record.KeyName,
		&record.Created,
		&record.Description,
		&isautogen)
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetSkAndProjectByAk(AccessKey string) (sk, pid string, err error) {
	var record AkSkRecord
	var isautogen bool
	err = Db.QueryRow("select * from aksk where accessKey=(?)", AccessKey).Scan(
		&record.AccessKey,
		&record.AccessSecret,
		&record.ProjectId,
		&record.AccountId,
		&record.KeyName,
		&record.Created,
		&record.Description,
		&isautogen)

	return record.AccessSecret, record.ProjectId, err
}

func RemoveAkSkRecord(AccessKey string) error {
	_, err := Db.Exec("delete from AkSk where accessKey=(?)", AccessKey)
	if err != nil {
		helper.Logger.Println(5, "Error RemoveAkSkRecord", AccessKey, err.Error())
	}
	return err
}

func ListAkSkRecordByProject(ProjectId string) ([]AccessKeyItem, error) {
	var items []AccessKeyItem
	var item AccessKeyItem
	rows, err := Db.Query("select * from AkSk where projectId=(?) and isAutoGenerated=false", ProjectId)
	if err != nil {
		helper.Logger.Println(5, "Error ListAkSkRecordByProject: ", err)
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		var record AkSkRecord
		var isautogen bool
		if err := rows.Scan(&record.AccessKey, &record.AccessSecret, &record.ProjectId, &record.AccountId, &record.KeyName, &record.Created, &record.Description, &isautogen); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}

		item.ProjectId = record.ProjectId
		item.AccessKey = record.AccessKey
		item.AccessSecret = record.AccessSecret
		item.Name = record.KeyName
		item.Status = "active" //fixme
		item.Updated = record.Created
		item.Created = record.Created //fixme
		item.Description = record.Description
		items = append(items, item)

	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return items, err
}

func GetKeysByAccessKeys(AccessKeys []string) ([]AccessKeyItem, error) {
	var items []AccessKeyItem
	var err error
	for _, key := range AccessKeys {
		var record AkSkRecord
		var item AccessKeyItem
		var isautogen bool
		err = Db.QueryRow("select * from AkSk where accessKey=(?)", key).Scan(
			&record.AccessKey,
			&record.AccessSecret,
			&record.ProjectId,
			&record.AccountId,
			&record.KeyName,
			&record.Created,
			&record.Description,
			&isautogen)
		if err != nil {
			helper.Logger.Println(5, "GetKeysByAccessKeys err: ", err)
			continue
		}
		item.ProjectId = record.ProjectId
		item.AccessKey = record.AccessKey
		item.AccessSecret = record.AccessSecret
		item.Name = record.KeyName
		item.Status = "active"
		item.Updated = record.Created
		item.Description = record.Description
		items = append(items, item)
	}
	return items, err
}

func GetKeysByAccount(accountid string) ([]AccessKeyItem, error) {
	var items []AccessKeyItem
	var err error
	var record AkSkRecord
	var item AccessKeyItem
	rows, err := Db.Query("select * from AkSk where accountid=(?)", accountid)
	if err != nil {
		helper.Logger.Println(5, "Error GetKeysByAccount: ", err)
		return items, err
	}

	defer rows.Close()

	var isautogen bool
	for rows.Next() {
		if err := rows.Scan(
			&record.AccessKey,
			&record.AccessSecret,
			&record.ProjectId,
			&record.AccountId,
			&record.KeyName,
			&record.Created,
			&record.Description,
			&isautogen); err != nil {
			helper.Logger.Println(5, "Row scan error: ", err)
			continue
		}

		item.ProjectId = record.ProjectId
		item.AccessKey = record.AccessKey
		item.AccessSecret = record.AccessSecret
		item.Name = record.KeyName
		item.Status = "active" //fixme
		item.Updated = record.Created
		item.Created = record.Created //fixme
		item.Description = record.Description
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		helper.Logger.Println(5, "Row error: ", err)
	}
	return items, err
}

func GetAutogenAkSkRecordByProject(ProjectId string) (AkSkRecord, error) {
	var record AkSkRecord
	var isautogen bool

	err := Db.QueryRow("select * from AkSk where projectId=(?) and isAutoGenerated=true", ProjectId).Scan(
		&record.AccessKey,
		&record.AccessSecret,
		&record.ProjectId,
		&record.AccountId,
		&record.KeyName,
		&record.Created,
		&record.Description,
		&isautogen)

	return record, err
}
