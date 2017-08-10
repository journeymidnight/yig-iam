package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api/datatype"
)

var Db *sql.DB

const (
	User_TYPE_ADMIN	= 0
	User_TYPE_ACCOUNT = 1
	User_TYPE_USER	= 2
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
	expired := time.Now().Add(time.Duration(helper.CONFIG.TokenExpire*1000000000)).Format(TimeFormat)
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

func InsertAkSkRecord(AccessKey string, SecretKey string, ProjectId string, AccountId string, KeyName string) error {
	created := time.Now().Format(TimeFormat)
	_, err := Db.Exec("insert into AkSk values( ?, ?, ?, ?, ?, ? )", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created)
	if err != nil {
		helper.Logger.Println(5, "Error InsertAkSkRecord", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created, err.Error())
	}
	return err
}

func IfAKExisted(AccessKey string) bool {
	var record AkSkRecord
	err := Db.QueryRow("select * from User where accessKey=(?)", AccessKey).Scan(
		&record.AccessKey,
		&record.AccessSecret,
		&record.ProjectId,
		&record.AccountId,
		&record.KeyName,
		&record.Created)
	if err != nil {
		return false
	} else {
		return true
	}
}

func RemoveAkSkRecord(AccessKey string, AccountId string) error {
	_, err := Db.Exec("delete from AkSk where accessKey=(?) and accountId=(?)", AccessKey, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error RemoveAkSkRecord", AccessKey, err.Error())
	}
	return err
}

func ListAkSkRecordByProject(ProjectId string, AccountId string) ([]AkSkRecord, error) {
	var records []AkSkRecord
	rows, err := Db.Query("select * from AkSk where projectId=(?) and accountId=(?)", ProjectId, AccountId)
	if err != nil {
		helper.Logger.Println(5, "Error ListAkSkRecordByProject: ", err)
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record AkSkRecord
		if err := rows.Scan(&record.AccessKey, &record.AccessSecret, &record.ProjectId, &record.AccountId, &record.KeyName, &record.Created); err != nil {
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

func GetKeysByAccessKeys(AccessKeys []string) ([]AccessKeyItem, error) {
	var items []AccessKeyItem
	var err error
	for _, key := range AccessKeys {
		var record AkSkRecord
		var item AccessKeyItem
		err = Db.QueryRow("select * from AkSk where accessKey=(?)", key).Scan(
			&record.AccessKey,
			&record.AccessSecret,
			&record.ProjectId,
			&record.AccountId,
			&record.KeyName,
			&record.Created)
		if err != nil {
			helper.Logger.Println(5, "GetKeysByAccessKeys err: ", err)
			continue
		}
		item.ProjectId = record.ProjectId
		item.AccessKey = record.AccessKey
		item.AccessSecret =  record.AccessSecret
		item.Name = record.KeyName
		item.Status = "active"
		item.Updated = record.Created
		items = append(items, item)
	}
	return items, err
}