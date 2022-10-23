package v1

import (
	"github.com/gin-gonic/gin"
	"mirror-api/model/kvDB/mirror"
	"mirror-api/util"
	"mirror-api/util/apiError"
	"net/http"
)

func MirrorsGET(w http.ResponseWriter, _ *http.Request, _ gin.Params) *apiError.Error {
	result := mirror.GetAll()

	return util.SendDataResponse(w, result, nil, http.StatusOK)
}
