package auth

// This plugin is based on Casbin: an authorization library that supports ACL, RBAC, ABAC
// View source at:
// https://github.com/casbin/casbin

import (
	"net/http"

	"context"
	"github.com/casbin/casbin"
	"github.com/journeymidnight/nier/src/api"
	. "github.com/journeymidnight/nier/src/api/datatype"
	"github.com/journeymidnight/nier/src/store"
	. "github.com/journeymidnight/nier/src/error"
	"github.com/journeymidnight/nier/src/helper"
	"github.com/urfave/negroni"
	"strings"
)

// Authz is a middleware that controls the access to the HTTP service, it is based
// on Casbin, which supports access control models like ACL, RBAC, ABAC.
// The plugin determines whether to allow a request based on (user, path, method).
// user: the authenticated user name.
// path: the URL for the requested resource.
// method: one of HTTP methods like GET, POST, PUT, DELETE.
//
// This middleware should be inserted fairly early in the middleware stack to
// protect subsequent layers. All the denied requests will not go further.
//
// It's notable that this middleware should be behind the authentication (e.g.,
// HTTP basic authentication, OAuth), so this plugin can get the logged-in user name
// to perform the authorization.
func Authorizer(e *casbin.Enforcer) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		isPublic := strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "env")
		if isPublic {
			next(w, r)
			return
		}
		isProxy := strings.Contains(r.URL.Path, "query") || strings.Contains(r.URL.Path, "lsblk")
		if isProxy {
			next(w, r)
			return
		}
		token := r.Header.Get("Token")
		if token == "" {
			api.WriteErrorResponse(w, r, ErrTokenEmpty)
			return
		}
		record, err := store.Db.GetTokenRecord(token)
		if err != nil {
			api.WriteErrorResponse(w, r, ErrTokenInvalid)
			return
		}
		//TODO: check if token has expired
		method := r.Method
		path := r.URL.Path
		helper.Logger.Infoln("authorize info", record.UserName, path, method)
		if strings.Contains(path, "catkeeper") {
			path = "/catkeeper"
		}
		authorised := e.Enforce(record.UserName, path, method)
		if authorised {
			ctx := context.WithValue(r.Context(), REQUEST_TOKEN_KEY, record)
			next(w, r.WithContext(ctx))
			return
		} else {
			api.WriteErrorResponse(w, r, ErrNotAuthorised)
			return
		}
	}
}

func GetUserName(r *http.Request) (string, string) {
	username, password, _ := r.BasicAuth()
	return username, password
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func CheckPermission(e *casbin.Enforcer, r *http.Request) bool {
	token := r.Header.Get("Token")
	if token == "" {
		return false
	}
	record, err := store.Db.GetTokenRecord(token)
	if err != nil {
		return false
	}
	//TODO: check if token has expired
	method := r.Method
	path := r.URL.Path
	if strings.Contains(path, "catkeeper") {
		path = "/catkeeper"
	}
	return e.Enforce(record.UserName, path, method)
}
