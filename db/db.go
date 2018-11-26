package db

import (
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/journeymidnight/yig-iam/helper"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	. "github.com/journeymidnight/yig-iam/error"
	"errors"
	"fmt"
)

var engine *xorm.Engine

const TimeFormat = "2006-01-02T15:04:05Z07:00"

func Db_Init() {
	var err error
	engine, err = xorm.NewEngine("mysql", helper.Config.UserDataSource)
	if err != nil {
		helper.Logger.Errorln(fmt.Sprintf("Connected to database %s failed", helper.Config.UserDataSource))
		panic(fmt.Sprintf("Connected to database %s failed", helper.Config.UserDataSource))
	}
	_, err = engine.Exec("CREATE DATABASE IF NOT EXISTS yig_iam")
	if err != nil {
		helper.Logger.Errorln("Create yig_iam failed")
		panic("fatal error: Create yig_iam failed")
	}
	engine, err = xorm.NewEngine("mysql", helper.Config.UserDataSource+"yig_iam")
	if err != nil {
		helper.Logger.Errorln(fmt.Sprintf("Connected to database %s failed", helper.Config.UserDataSource+"yig_iam"))
		panic(fmt.Sprintf("Connected to database %s failed", helper.Config.UserDataSource+"yig_iam"))
	}
	exist, err := engine.IsTableExist(&User{})
	if err != nil {
		panic("check table failed")
	} else {
		if exist == false {
			if engine.CreateTables(&User{}) != nil {
				panic("create table failed")
			}
		}
	}
	exist, err = engine.IsTableExist(&Token{})
	if err != nil {
		panic("check table failed")
	} else {
		if exist == false {
			if engine.CreateTables(&Token{}) != nil {
				panic("create table failed")
			}
		}
	}
	exist, err = engine.IsTableExist(&Project{})
	if err != nil {
		panic("check table failed")
	} else {
		if exist == false {
			if engine.CreateTables(&Project{}) != nil {
				panic("create table failed")
			}
		}
	}
	exist, err = engine.IsTableExist(&UserProject{})
	if err != nil {
		panic("check table failed")
	} else {
		if exist == false {
			if engine.CreateTables(&UserProject{}) != nil {
				panic("create table failed")
			}
		}
	}
	var root User
	root.UserName = "root"
	root.Password = "root"
	root.UserId = helper.GenerateUserId()
	root.Type = ROLE_ROOT
	root.DisplayName = "default_root"
	root.AccountId = root.UserId
	root.Status = USER_STATUS_ACTIVE
	engine.Insert(&root)

	var admin User
	admin.UserName = "admin"
	admin.Password = "admin"
	admin.UserId = helper.GenerateUserId()
	admin.Type = ROLE_ACCOUNT
	admin.DisplayName = "default_account"
	admin.AccountId = admin.UserId
	admin.Status = USER_STATUS_ACTIVE
	engine.Insert(&admin)

	engine.ShowSQL(true)
}

func Engine() *xorm.Engine {
	if engine == nil {
		Db_Init()
		engine.ShowSQL(true)
	}
	return engine
}

func DeleteAccount(userName string) error {
	//_, err := Db.Exec("delete from User where accountId=(?) and type='ACCOUNT'", AccountId)
	var user User
	user.UserName = userName
	user.Type = ROLE_ACCOUNT
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe account", userName, err.Error())
		return err
	}
	if has == false {
		return ErrDbRecordNotFound
	}

	session := Engine().NewSession()
	defer session.Close()
	err = session.Begin()

	_, err = session.Delete(&User{AccountId:user.AccountId})
	if err != nil {
		helper.Logger.Errorln("Error delete all users for account:", userName, err.Error())
		session.Rollback()
		return err
	}

	_, err = session.Delete(&UserProject{AccountId:user.AccountId})
	if err != nil {
		helper.Logger.Errorln("Error delete all user-project for account:", userName, err.Error())
		session.Rollback()
		return err
	}

	_, err = Engine().Update(&Project{Status:PROJECT_STATUS_DELETED}, &Project{AccountId:user.AccountId})
	if err != nil {
		helper.Logger.Errorln("Error mark all projects deleted for account:", userName, err.Error())
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DescribeAccount(userName string) (user User, err error) {
	user.UserName = userName
	user.Type = ROLE_ACCOUNT
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe account", userName, err.Error())
		return user, err
	}
	if has {
		return user, nil
	} else {
		return user, ErrDbRecordNotFound
	}
}

func DeactivateAccount(userName string) (error) {
	var user User
	user.UserName = userName
	user.Type = ROLE_ACCOUNT
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe account", userName, err.Error())
		return err
	}
	if has == false {
		return ErrDbRecordNotFound
	}

	_, err = Engine().Update(&UserProject{Status:KEY_STATUS_DISABLE}, &UserProject{AccountId:user.AccountId})
	if err != nil {
		helper.Logger.Errorln("Error deactivate account", userName, err.Error())
		return err
	}

	return nil
}

func ActivateAccount(userName string) (error) {
	var user User
	user.UserName = userName
	user.Type = ROLE_ACCOUNT
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe account", userName, err.Error())
		return err
	}
	if has == false {
		return ErrDbRecordNotFound
	}

	_, err = Engine().Update(&UserProject{Status:KEY_STATUS_ENABLE}, &UserProject{AccountId:user.AccountId})
	//_, err := Db.Exec("update User set status='inactive' where accountId=(?) and type='ACCOUNT'", AccountId)
	if err != nil {
		helper.Logger.Errorln("Error activate account", userName, err.Error())
		return err
	}

	return nil
}

func ListAccounts() ([]User, error) {
	users := make([]User, 0)
	err := Engine().Where("Type = ?", ROLE_ACCOUNT).Find(&users)
	if err != nil {
		helper.Logger.Errorln("Error querying idle executors: ", err)
		return nil, err
	}
	return users, nil
}

func CreateUser(userName string, password string, accountType string,
			email string, displayName string, accountId string, createProject bool) error {
	session := Engine().NewSession()
	defer session.Close()
	err := session.Begin()

	var user User
	user.UserName = userName
	user.Password = password


	if displayName == "" {
		user.DisplayName = userName
	}else {
		user.DisplayName = displayName
	}
	user.Type = accountType
	user.Email = email
	user.UserId = "u-" + string(helper.GenerateRandomId())
	if accountType == ROLE_ACCOUNT {
		user.AccountId = user.UserId
	} else {
		user.AccountId = accountId
	}
	user.Status = USER_STATUS_ACTIVE

	_, err = session.Insert(&user)
	if err != nil {
		helper.Logger.Errorln("Error create user", userName, accountType, accountId, err.Error())
		session.Rollback()
		return err
	}

	if createProject == true {
		var p Project
		projectId := "p-" + string(helper.GenerateRandomId())
		p.ProjectId = projectId
		p.ProjectName = "Private"
		p.ProjectType = PRIVATE_PROJECT
		p.AccountId = accountId
		p.OwnerId = user.UserId
		p.Status = PROJECT_STATUS_ACTIVE
		p.Description = fmt.Sprintf("own by %s", userName)
		_, err = session.Insert(&p)
		if err != nil {
			helper.Logger.Errorln("Error create private project", projectId, userName, err.Error())
			session.Rollback()
			return err
		}

		var up UserProject
		ak, sk := helper.GenerateKey()
		up.AccessKey = string(ak)
		up.AccessSecret = string(sk)
		up.UserId = user.UserId
		up.ProjectId = p.ProjectId
		up.AccountId = accountId
		up.Status = KEY_STATUS_ENABLE
		up.Acl = ACL_RW
		_, err = session.Insert(&up)
		if err != nil {
			helper.Logger.Errorln("Error create user-project", projectId, up.UserId, accountId, up.Acl, err.Error())
			session.Rollback()
			return err
		}
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	return nil
}

func RemoveUser(userName string, accountId string) error {
	var user User
	user.UserName = userName
	user.AccountId = accountId
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe user", userName, err.Error())
		return err
	}
	if has == false {
		return errors.New("user not existed")
	}

	session := Engine().NewSession()
	defer session.Close()
	err = session.Begin()

	_, err = session.Delete(&user)
	if err != nil{
		helper.Logger.Errorln("Error remove user:", userName, err.Error())
		session.Rollback()
		return err
	}

	var up UserProject
	up.UserId = user.UserId
	_, err = session.Delete(&up)
	if err != nil{
		helper.Logger.Errorln("Error remove user-project:", userName, err.Error())
		session.Rollback()
		return err
	}

	var p Project
	p.Status = PROJECT_STATUS_DELETED
	_, err = session.Update(&p, &Project{OwnerId:user.UserId})
	if err != nil{
		helper.Logger.Errorln("Error remove user-project:", userName, err.Error())
		session.Rollback()
		return err
	}
	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user User) error {
	_, err := Engine().Update(&user, &User{UserId:user.UserId})
	if err != nil {
		helper.Logger.Errorln("Error mark all projects deleted for account:", user.UserId, err.Error())
		return err
	}
	return nil
}

func DescribeUser(userName string, accountId string) (user User, err error) {
	user.UserName = userName
	user.AccountId = accountId
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error describe user", userName, err.Error())
		return user, err
	}
	if has {
		return user, nil
	} else {
		return user, ErrDbRecordNotFound
	}
}

func ValidateEmailAndPassword(email, password string) (user User, err error) {
	user.Email = email
	user.Password = password
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error validate user", email, password, err.Error())
		return user, err
	}
	if has {
		return user, nil
	} else {
		return user, ErrDbRecordNotFound
	}
}

func ValidateUserAndPassword(userName string, password string) (user User, err error) {
	user.UserName = userName
	user.Password = password
	has, err := Engine().Get(&user)
	if err != nil {
		helper.Logger.Errorln("Error validate user", userName, password, err.Error())
		return user, err
	}
	if has {
		return user, nil
	} else {
		return user, ErrDbRecordNotFound
	}
}

func ListUsers(accountId string) (users []User, err error) {
	err = Engine().Where("accountId = ? AND type = ?", accountId, ROLE_USER).Find(&users)
	if err != nil {
		helper.Logger.Errorln("Error list users", accountId, ROLE_USER, err.Error())
	}
	return
}

func CreateProject(projectName, projectType, accountId, description string) error {
	projectId := "p-" + string(helper.GenerateRandomId())
	session := Engine().NewSession()
	defer session.Close()
	err := session.Begin()

	var p Project
	p.ProjectId = projectId
	p.ProjectName = projectName
	p.ProjectType = projectType
	p.AccountId = accountId
	p.Description = description
	p.OwnerId = accountId
	p.Status = PROJECT_STATUS_ACTIVE
	_, err = session.Insert(&p)
	if err != nil {
		helper.Logger.Errorln("Error create project", projectId, projectName, accountId, projectType, err.Error())
		session.Rollback()
		return err
	}
	var up UserProject
	ak, sk := helper.GenerateKey()
	up.AccessKey = string(ak)
	up.AccessSecret = string(sk)
	up.UserId = accountId
	up.ProjectId = projectId
	up.Acl = ACL_RW
	up.Status = KEY_STATUS_ENABLE
	up.AccountId = accountId
	_, err = session.Insert(&up)
	if err != nil {
		helper.Logger.Errorln("Error create user-project", projectId, projectName, accountId, projectType, err.Error())
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	return err
}

func RemoveProject(projectId string, accountId string) error {
	target := &Project{ProjectId:projectId}
	has, err := Engine().Get(target)
	if err != nil {
		return err
	}

	if has == false {
		return ErrDbRecordNotFound
	}

	if target.AccountId != accountId {
		return ErrNotAuthorised
	}

	session := Engine().NewSession()
	defer session.Close()
	err = session.Begin()
	p := Project{Status:PROJECT_STATUS_DELETED}
	_, err = session.Update(&p, &Project{ProjectId:projectId})
	if err != nil {
		session.Rollback()
		return err
	}

	up := UserProject{ProjectId:projectId}
	_, err = session.Delete(&up)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DescribeProject(projectId string, accountId string) (p Project, err error) {
	p.ProjectId = projectId
	p.AccountId = accountId
	has, err := Engine().Get(&p)
	if err != nil {
		helper.Logger.Errorln("Error describe project:", projectId, err.Error())
		return p, err
	}
	if has {
		return p, nil
	} else {
		return p, ErrDbRecordNotFound
	}
}

func ListProjects(userId string) ([]ListProjectResp, error) {
	resp := make([]ListProjectResp, 0)
	Engine().Join("INNER", "project", "project.projectId = user_project.projectId").Where("user_project.userId=? and project.status <>?",userId,PROJECT_STATUS_DELETED).Find(&resp)
	return resp, nil
}

func ListProjectsByAccount(accountId string) (projects []Project, err error) {
	err = Engine().Where("accountId = ?", accountId).Find(&projects)
	if err != nil {
		helper.Logger.Errorln("Error list projects", accountId, err.Error())
		return projects, err
	}
	return
}

func ListProjectsByOwner(ownerId string) (projects []Project, err error) {
	err = Engine().Where("ownerId = ?", ownerId).Find(&projects)
	if err != nil {
		helper.Logger.Errorln("Error list projects", ownerId, err.Error())
		return projects, err
	}
	return
}

func LinkUserWithProject(projectId string, userId string, acl string, accountId string) error {
	var up UserProject
	up.ProjectId = projectId
	up.UserId = userId
	up.AccessKey = string(helper.GenerateRandomIdByLength(20))
	up.AccessSecret = string(helper.GenerateRandomIdByLength(40))
	up.Acl = acl
	up.Status = KEY_STATUS_ENABLE
	up.AccountId = accountId
	_, err := Engine().Insert(&up)
	if err != nil{
		helper.Logger.Errorln("Error link user with project", projectId, userId, acl, err.Error())
	}
	return err
}

func UnlinkUserWithProject(projectId string, userId string) error {
	var up UserProject
	up.ProjectId = projectId
	up.UserId = userId
	affected, err := Engine().Delete(&up)
	if err != nil{
		helper.Logger.Errorln("Error unlink project user:", projectId, userId, err.Error())
		return err
	}
	if affected > 0 {
		return nil
	} else {
		helper.Logger.Println("Project User record not existed:", projectId, userId)
		return ErrDbRecordNotFound
	}
	return err
}

//func ListUserProjectLinksByUser(userId string) (ups []UserProject, err error) {
//	err = Engine().Where("userId = ?", userId).Find(&ups)
//	if err != nil {
//		helper.Logger.Errorln("Error list projects", userId, err.Error())
//	}
//	return
//}

func ListUsersByProject(projectId string) (ups []ListUserResp, err error) {
	resp := make([]ListUserResp, 0)
	Engine().Join("INNER", "user", "user.userId = user_project.userId").Where("user_project.projectId=?",projectId).Find(&resp)
	return resp, nil
}

func CheckAccountIdExist(accountId string) (exist bool, err error) {
	var user User
	user.AccountId = accountId
	user.Type = ROLE_ACCOUNT
	exist, err = Engine().Exist(&user)
	return
}

func CheckUserExist(userId string) (exist bool, err error) {
	var user User
	user.UserId = userId
	exist, err = Engine().Exist(&user)
	return

}

func CreateToken(token string, userId string, userName string, accountId string, userType string) error {
	var t Token
	t.Token = token
	t.UserName = userName
	t.UserId = userId
	t.Type = userType
	t.AccountId = accountId
	_, err := Engine().Insert(&t)
	if err != nil{
		helper.Logger.Errorln("Error create token", token, userId, userName, err.Error())
	}
	return err
}

func GetToken(token string) (t Token, err error) {
	t.Token = token
	has, err := Engine().Get(&t)
	if err != nil {
		helper.Logger.Errorln("Error describe token:", token, err.Error())
		return t, err
	}
	if has {
		return t, nil
	} else {
		return t, ErrDbRecordNotFound
	}
}

func RemoveToken(token string) error {
	var t Token
	t.Token = token
	affected, err := Engine().Delete(&t)
	if err != nil{
		helper.Logger.Errorln("Error remove token:", token, err.Error())
		return err
	}
	if affected > 0 {
		return nil
	} else {
		helper.Logger.Println("token not existed:", token)
		return ErrDbRecordNotFound
	}
}

func SearchExistedToken(userName string) (token Token, err error) {
	token.UserName = userName
	has, err := Engine().Get(&token)
	if err != nil {
		helper.Logger.Errorln("Error describe token:", token, err.Error())
		return token, err
	}
	if has {
		return token, nil
	} else {
		return token, ErrDbRecordNotFound
	}
}

func ListExpiredTokens() (tokens []Token, err error) {
	now := time.Now()
	expired := now.Add(-time.Duration(helper.Config.TokenExpire * 1000000000))


	err = Engine().Iterate(new(Token), func(i int, bean interface{})error {
		token := bean.(*Token)
		if expired.Sub(time.Time(token.CreatedAt)) >0 {
			tokens = append(tokens, *token)
		}
		return nil
	})
	if err != nil {
		helper.Logger.Errorln("Error list expired tokens", err.Error())
	}
	return
}

//func InsertAkSkRecord(AccessKey string, SecretKey string, ProjectId string, AccountId string, KeyName string, Description string) error {
//	created := time.Now().Format(TimeFormat)
//	_, err := Db.Exec("insert into AkSk values( ?, ?, ?, ?, ?, ?, ?)", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created, Description)
//	if err != nil {
//		helper.Logger.Println(5, "Error InsertAkSkRecord ", AccessKey, SecretKey, ProjectId, AccountId, KeyName, created, Description,err.Error())
//	}
//	return err
//}
//
//func IfAKExisted(AccessKey string) bool {
//	var record AkSkRecord
//	err := Db.QueryRow("select * from User where accessKey=(?)", AccessKey).Scan(
//		&record.AccessKey,
//		&record.AccessSecret,
//		&record.ProjectId,
//		&record.AccountId,
//		&record.KeyName,
//		&record.Created,
//		&record.Description)
//	if err != nil {
//		return false
//	} else {
//		return true
//	}
//}

//func RemoveAkSkRecord(AccessKey string, AccountId string) error {
//	_, err := Db.Exec("delete from AkSk where accessKey=(?) and accountId=(?)", AccessKey, AccountId)
//	if err != nil {
//		helper.Logger.Println(5, "Error RemoveAkSkRecord", AccessKey, err.Error())
//	}
//	return err
//}

//func ListAkSkRecordByProject(ProjectId string, AccountId string) ([]AkSkRecord, error) {
//	var records []AkSkRecord
//	rows, err := Db.Query("select * from AkSk where projectId=(?) and accountId=(?)", ProjectId, AccountId)
//	if err != nil {
//		helper.Logger.Println(5, "Error ListAkSkRecordByProject: ", err)
//		return records, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var record AkSkRecord
//		if err := rows.Scan(&record.AccessKey, &record.AccessSecret, &record.ProjectId, &record.AccountId, &record.KeyName, &record.Created, &record.Description); err != nil {
//			helper.Logger.Println(5, "Row scan error: ", err)
//			continue
//		}
//		records = append(records, record)
//	}
//	if err := rows.Err(); err != nil {
//		helper.Logger.Println(5, "Row error: ", err)
//	}
//	return records, err
//}
//
//func ListKeyRecordsByProjects(ProjectIds []string, AccountId string) ([]AkSkRecord, error) {
//	var records []AkSkRecord
//	rows, err := Db.Query("select * from AkSk where projectId=(?) and accountId=(?)", ProjectId, AccountId)
//	if err != nil {
//		helper.Logger.Println(5, "Error ListAkSkRecordByProject: ", err)
//		return records, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var record AkSkRecord
//		if err := rows.Scan(&record.AccessKey, &record.AccessSecret, &record.ProjectId, &record.AccountId, &record.KeyName, &record.Created, &record.Description); err != nil {
//			helper.Logger.Println(5, "Row scan error: ", err)
//			continue
//		}
//		records = append(records, record)
//	}
//	if err := rows.Err(); err != nil {
//		helper.Logger.Println(5, "Row error: ", err)
//	}
//	return records, err
//}

func GetKeyItemsByAccessKeys(AccessKeys []string) ([]AccessKeyItem, error) {
	var err error

	items := make([]AccessKeyItem, 0)
	for _, key := range AccessKeys {
		var up UserProject
		up.AccessKey = key
		has, err := Engine().Get(&up)
		if err != nil {
			helper.Logger.Errorln("Error GetKeyItemsByAccessKeys by accessKey", key, err.Error())
			return items, err
		}
		if has {
			var item AccessKeyItem
			item.ProjectId = up.ProjectId
			item.Name = up.ProjectId
			item.AccessKey = up.AccessKey
			item.AccessSecret = up.AccessSecret
			item.Acl = up.Acl
			item.Status = up.Status
			item.Updated = time.Time(up.UpdatedAt).Format("2006-01-02 15:04:05")
			items = append(items, item)
		}
	}
	return items, err
}

func GetKeyItemByAccessKey(AccessKey string) ([]AccessKeyItem, error) {
	var err error
	var up UserProject
	items := make([]AccessKeyItem, 0)

	up.AccessKey = AccessKey
	has, err := Engine().Get(&up)
	if err != nil {
		helper.Logger.Errorln("Error GetKeyItemByAccessKeyt by accessKey", AccessKey, err.Error())
		return items, err
	}
	if has {
		var item AccessKeyItem
		item.ProjectId = up.ProjectId
		item.Name = up.ProjectId
		item.AccessKey = up.AccessKey
		item.AccessSecret = up.AccessSecret
		item.Acl = up.Acl
		item.Status = up.Status
		item.Updated = time.Time(up.UpdatedAt).Format("2006-01-02 15:04:05")
		items = append(items, item)
	}
	return items, nil
}

func GetKeyItemsByProject(projectId string) ([]AccessKeyItem, error) {
	var err error
	ups := make([]UserProject, 0)
	items := make([]AccessKeyItem, 0)
	err = Engine().Where("user_project.projectId=?", projectId).Find(&ups)
	if err != nil {
		helper.Logger.Println("Error GetKeyItemsByProject: ", err)
		return items, err
	}
	for _, up := range ups {
		var item AccessKeyItem
		item.ProjectId = up.ProjectId
		item.Name = up.ProjectId
		item.AccessKey = up.AccessKey
		item.AccessSecret = up.AccessSecret
		item.Acl = up.Acl
		item.Status = up.Status
		item.Updated = time.Time(up.UpdatedAt).Format("2006-01-02 15:04:05")
		items = append(items, item)
	}
	return items, err
}

func GetKeyItemsByUserId(useId string) ([]AccessKeyItem, error) {
	var err error
	ups := make([]UserProject, 0)
	items := make([]AccessKeyItem, 0)
	err = Engine().Where("user_project.userId=?", useId).Find(&ups)
	if err != nil {
		helper.Logger.Println("Error GetKeyItemsByProject: ", err)
		return items, err
	}
	for _, up := range ups {
		var item AccessKeyItem
		item.ProjectId = up.ProjectId
		item.Name = up.ProjectId
		item.AccessKey = up.AccessKey
		item.AccessSecret = up.AccessSecret
		item.Acl = up.Acl
		item.Status = up.Status
		item.Updated = time.Time(up.UpdatedAt).Format("2006-01-02 15:04:05")
		items = append(items, item)
	}
	return items, err
}