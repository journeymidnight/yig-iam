package api

import (
	"net/http"
	. "github.com/journeymidnight/yig-iam/error"
	"github.com/journeymidnight/yig-iam/helper"
	"encoding/json"
)

// APIErrorResponse - error response format
type ApiErrorResponse struct {
	ErrorCode string
	Message   string
}

type ApiUserLoginResponse struct {
	Token string
	Type string
}

func EncodeResponse(response interface{}) []byte {
	res, _ := json.Marshal(response)
	return res
}

// WriteSuccessResponse write success headers and response if any.
func WriteSuccessResponse(w http.ResponseWriter, response []byte) {
	if response == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Write(response)
	w.(http.Flusher).Flush()
}

// writeSuccessNoContent write success headers with http status 204
func WriteSuccessNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// writeErrorResponse write error headers
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteErrorResponseHeaders(w, err)
	WriteErrorResponseNoHeader(w, r, err, r.URL.Path)
}

func WriteErrorResponseWithResource(w http.ResponseWriter, r *http.Request, err error, resource string) {
	WriteErrorResponseHeaders(w, err)
	WriteErrorResponseNoHeader(w, r, err, resource)
}

func WriteErrorResponseHeaders(w http.ResponseWriter, err error) {
	var status int
	apiErrorCode, ok := err.(ApiError)
	if ok {
		status = apiErrorCode.HttpStatusCode()
	} else {
		status = http.StatusInternalServerError
	}
	helper.Logger.Infoln("Response status code:", status)
	w.WriteHeader(status)
}

func WriteErrorResponseNoHeader(w http.ResponseWriter, req *http.Request, err error, resource string) {
	// HEAD should have no body, do not attempt to write to it
	if req.Method == "HEAD" {
		return
	}

	// Generate error response.
	errorResponse := ApiErrorResponse{}
	apiErrorCode, ok := err.(ApiError)
	if ok {
		errorResponse.ErrorCode = apiErrorCode.ErrorCode()
		errorResponse.Message = apiErrorCode.Description()
	} else {
		errorResponse.ErrorCode = "InternalError"
		errorResponse.Message = err.Error()
	}

	encodedErrorResponse := EncodeResponse(errorResponse)
	w.Write(encodedErrorResponse)
	w.(http.Flusher).Flush()
}