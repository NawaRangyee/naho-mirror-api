package util

import (
	json "github.com/json-iterator/go"
	"mirror-api/model/_util/pageInfo"
	"mirror-api/util/apiError"
	"net/http"
)

type Response struct {
	Code       int                `json:"code"`
	Message    string             `json:"message,omitempty"`
	Success    bool               `json:"success"`
	PagingInfo *pageInfo.Response `json:"paging_info,omitempty"`
	Data       *interface{}       `json:"data,omitempty"`
}

type MResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type DResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func SendJSONResponse(w http.ResponseWriter, body interface{}) *apiError.Error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return apiError.InternalServerError(err)
	}

	return nil
}

func SendDataResponse(w http.ResponseWriter, data interface{}, pi *pageInfo.Response, code int) *apiError.Error {
	resp := &Response{Code: code, Message: "", Success: code >= 200 && code < 300, PagingInfo: pi}
	if data != nil {
		resp.Data = &data
	}
	if code >= 100 && code < 600 {
		w.WriteHeader(code)
	}

	return SendJSONResponse(w, resp)
}

func SendMessageResponse(w http.ResponseWriter, msg string, code int) *apiError.Error {
	resp := &Response{Code: code, Message: msg, Success: code >= 200 && code < 300}
	if code >= 100 && code < 600 {
		w.WriteHeader(code)
	}

	return SendJSONResponse(w, resp)
}

func SendNoContent(w http.ResponseWriter) *apiError.Error {
	w.WriteHeader(http.StatusNoContent)

	return nil
}
