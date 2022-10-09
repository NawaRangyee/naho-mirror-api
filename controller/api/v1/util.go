package v1

import (
	"github.com/gin-gonic/gin"
	"mirror-api/config"
	"mirror-api/util"
	"mirror-api/util/apiError"
	"net/http"
	"strings"
	"time"
)

type InfoResp struct {
	ServerTimeUnix int64  `json:"server_time"`
	TimeZone       string `json:"time_zone"`
	Hostname       string `json:"hostname"`
}

const (
	queryAllOK = "query_all_ok"
)

func CoffeeGET(w http.ResponseWriter, _ *http.Request, _ gin.Params) *apiError.Error {
	return util.SendMessageResponse(w, "I reject to boil your teapot.", http.StatusTeapot)
}

func Ping(w http.ResponseWriter, _ *http.Request, _ gin.Params) *apiError.Error {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return apiError.InternalServerError(err)
	}

	return nil
}

func Info(w http.ResponseWriter, _ *http.Request, _ gin.Params) *apiError.Error {
	return util.SendDataResponse(w, getInfo(), nil, http.StatusOK)
}

func getContentType(r *http.Request) string {
	return strings.ToLower(r.Header.Get("Content-Type"))
}

func getInfo() *InfoResp {
	return &InfoResp{ServerTimeUnix: time.Now().Unix(), TimeZone: config.Loc.String(), Hostname: config.GetHostname()}
}
