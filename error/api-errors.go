package error

import (
	"net/http"
)

type ApiError interface {
	error
	ErrorCode() string
	Description() string
	HttpStatusCode() int
}

type ApiErrorStruct struct {
	ErrorCode      string
	Description    string
	HttpStatusCode int
}

// APIErrorCode type of error status.
type ApiErrorCode int

// Error codes, non exhaustive list - http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
const (
	ErrUserExised ApiErrorCode = iota
	ErrTokenEmpty
	ErrTokenInvalid
	ErrTokenExpired
	ErrNotAuthorised
	ErrJsonDecodeFailed
	ErrUserOrPasswordInvalid
	ErrAccountDisabled
	ErrFailedAddUser
	ErrFailedCreateProject
	ErrInvalidUserType
	ErrInvalidUserPassword
	ErrDuplicateAddUser
	ErrInternalError
	ErrDbOperateFailed
	ErrDbRecordNotFound
	ErrInvalidParameters
)

// error code to APIError structure, these fields carry respective
// descriptions for all the error responses.
var ErrorCodeResponse = map[ApiErrorCode]ApiErrorStruct{
	//ErrInvalidCopyDest: {
	//	AwsErrorCode:   "InvalidRequest",
	//	Description:    "This copy request is illegal because it is trying to copy an object to itself.",
	//	HttpStatusCode: http.StatusBadRequest,
	//},
	//
	//ErrSignatureVersionNotSupported: {
	//	AwsErrorCode:   "AccessDenied",
	//	Description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
	//	HttpStatusCode: http.StatusForbidden,
	//},
	ErrUserExised: {
		ErrorCode:      "User already existed",
		Description:    "The user you tried to create is existed",
		HttpStatusCode: http.StatusConflict,
	},
	ErrTokenEmpty: {
		ErrorCode:      "Token empty",
		Description:    "Should login first",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrTokenInvalid: {
		ErrorCode:      "Token invalid",
		Description:    "Should login again",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrTokenExpired: {
		ErrorCode:      "Token expired",
		Description:    "Should login again",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrNotAuthorised: {
		ErrorCode:      "Not Authorised",
		Description:    "Do not have authority to perform this action",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrJsonDecodeFailed: {
		ErrorCode:      "Body decode failed",
		Description:    "received bad request body",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrUserOrPasswordInvalid: {
		ErrorCode:      "User or Password Invalid",
		Description:    "User or Password Invalid",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrAccountDisabled: {
		ErrorCode:      "Account disabled",
		Description:    "Account disabled",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrFailedAddUser: {
		ErrorCode:      "Failed add user",
		Description:    "Failed add user",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrFailedCreateProject: {
		ErrorCode:      "Failed create project",
		Description:    "Failed create project",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrInvalidUserType: {
		ErrorCode:      "Invalid User Type",
		Description:    "User Type Could Only Be [ADMIN/USER]",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrInvalidUserPassword: {
		ErrorCode:      "Invalid User Password",
		Description:    "Invalid User Password]",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrDuplicateAddUser: {
		ErrorCode:      "user already exist",
		Description:    "user already exist",
		HttpStatusCode: http.StatusConflict,
	},
	ErrInternalError: {
		ErrorCode:      "Internal Error",
		Description:    "Internal Error",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrDbOperateFailed: {
		ErrorCode:      "database operate failed",
		Description:    "database operate failed",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrDbRecordNotFound: {
		ErrorCode:      "database record not existed",
		Description:    "database record not existed",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrInvalidParameters: {
		ErrorCode:      "parameters in body not complete",
		Description:    "parameters in body not complete",
		HttpStatusCode: http.StatusBadRequest,
	},


	//ErrBucketAccessForbidden: {
	//	AwsErrorCode:   "AccessDenied",
	//	Description:    "You have no access to this bucket.",
	//	HttpStatusCode: http.StatusForbidden,
	//},
	//
	//ErrNoSuchBucketLc: {
	//	AwsErrorCode:   "NoSuchBucketLc",
	//	Description:    "The specified bucket does not have LifeCycle configured.",
	//	HttpStatusCode: http.StatusNotFound,
	//},
}

func (e ApiErrorCode) ErrorCode() string {
	awsError, ok := ErrorCodeResponse[e]
	if !ok {
		return "InternalError"
	}
	return awsError.ErrorCode
}

func (e ApiErrorCode) Description() string {
	nierError, ok := ErrorCodeResponse[e]
	if !ok {
		return "We encountered an internal error, please try again."
	}
	return nierError.Description
}

func (e ApiErrorCode) Error() string {
	return e.Description()
}

func (e ApiErrorCode) HttpStatusCode() int {
	nierError, ok := ErrorCodeResponse[e]
	if !ok {
		return http.StatusInternalServerError
	}
	return nierError.HttpStatusCode
}

