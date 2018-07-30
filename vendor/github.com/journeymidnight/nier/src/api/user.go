package api

import (
	"net/http"
	"io/ioutil"
	"github.com/journeymidnight/nier/src/store"
	. "github.com/journeymidnight/nier/src/api/datatype"
	"github.com/journeymidnight/nier/src/helper"
	"encoding/json"
	"github.com/google/uuid"
	. "github.com/journeymidnight/nier/src/error"
	"github.com/go-sql-driver/mysql"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	switch req.Type {
	case ROLE_ADMIN:
	case ROLE_USER:
	default:
		req.Type = ROLE_USER
	}

	err = store.Db.InsertUserRecord(req.Name, req.Password, req.Type)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok {
			helper.Logger.Infoln("failed CreateUser for query:", req)
			WriteErrorResponse(w, r, err)
			return
		}
		if me.Number == 1062 {
			WriteErrorResponse(w, r, ErrDuplicateAddUser)
			return
		}
		helper.Logger.Infoln("failed CreateUser for query:", req)
		WriteErrorResponse(w, r, ErrFailedAddUser)
		return
	}
	helper.Enforcer.AddRoleForUser(req.Name, req.Type)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		helper.Logger.Errorln("decode error:", err.Error())
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	record, err := store.Db.ValidUserAndPassword(req.Name, req.Password)
	if err != nil {
		WriteErrorResponse(w, r, ErrUserOrPasswordInvalid)
		return
	}
	token, err := store.Db.SearchExistedToken(record.UserName)
	if err == nil {
		WriteSuccessResponse(w, EncodeResponse(ApiUserLoginResponse{token.Token, token.Type}))
		return
	}
	uuid := uuid.New()
	err = store.Db.InsertTokenRecord(uuid.String(), record.UserName, record.Type)
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(ApiUserLoginResponse{uuid.String(), record.Type}))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if req.Name == "admin" {
		WriteErrorResponse(w, r, ErrNotAuthorised)
		return
	}
	err = store.Db.RemoveUserRecord(req.Name)
	if err != nil {
		helper.Logger.Infoln("failed DeleteUser for query:", req)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Enforcer.DeleteRolesForUser(req.Name)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &QueryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}
	if req.Password == "" {
		WriteErrorResponse(w, r, ErrInvalidUserPassword)
		return
	}
	switch req.Type {
	case ROLE_ADMIN:
	case ROLE_USER:
	default:
		if req.Name == "admin" {
			break
		}
		WriteErrorResponse(w, r, ErrInvalidUserType)
		return
	}
	err = store.Db.ModifyUserRecord(req.Name, req.Password, req.Type)
	if err != nil {
		helper.Logger.Infoln("failed DeleteUser for query:", req)
		WriteErrorResponse(w, r, err)
		return
	}
	helper.Enforcer.DeleteRolesForUser(req.Name)
	helper.Enforcer.AddRoleForUser(req.Name, req.Type)
	helper.Enforcer.SavePolicy()
	WriteSuccessResponse(w, nil)
}

func ListUser(w http.ResponseWriter, r *http.Request) {
	records, err := store.Db.ListUserRecords()
	if err != nil {
		helper.Logger.Errorln("failed list user records:")
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(records))
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	WriteSuccessResponse(w, EncodeResponse(Roles{helper.Enforcer.GetAllRoles()}))
}
