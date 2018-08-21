package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/google/uuid"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	. "github.com/journeymidnight/yig-iam/error"
	"github.com/journeymidnight/yig-iam/db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := &QueryRequest{}
	err := json.Unmarshal(body, query)
	if err != nil {
		WriteErrorResponse(w, r, ErrJsonDecodeFailed)
		return
	}

	user, err := db.ValidateUserAndPassword(query.UserName, query.Password)
	if err != nil {
		//retry auth by email
		if query.Email != "" {
			user, err = db.ValidateEmailAndPassword(query.Email, query.Password)
			if err != nil {
				WriteErrorResponse(w, r, ErrUserOrPasswordInvalid)
				return
			}
		}
	}

	if user.Status == USER_STATUS_INACTIVE {
		WriteErrorResponse(w, r, ErrAccountDisabled)
		return
	}

	uuid := uuid.New()
	err = db.CreateToken(uuid.String(), user.UserId, user.UserName, user.AccountId, user.Type)
	if err != nil {
		WriteErrorResponse(w, r, err)
		return
	}
	var resp LoginResponse
	resp.Token = uuid.String()
	resp.Type = user.Type
	resp.AccountId = user.AccountId
	resp.UserId = user.UserId

	WriteSuccessResponse(w, EncodeResponse(resp))
	return

}
