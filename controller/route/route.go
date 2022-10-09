package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "mirror-api/controller/api/v1"
	"mirror-api/util"
	"net/http"
)

func Load() http.Handler {
	return middleware(routes())
}

func getEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/ping", "/coffee", "/info"}}))
	router.Use(gin.Recovery())

	return router
}

func routes() *gin.Engine {
	// Force log's color
	gin.ForceConsoleColor()

	router := getEngine()

	// MARKS : API v1

	// Util
	router.GET("/coffee", handle(v1.CoffeeGET))
	router.GET("/info", handle(v1.Info))

	// Ping-Pong
	router.GET("/ping", handle(v1.Ping))

	return router
}

// handle wraps out api handler to router
func handle(handler util.APIHandler) gin.HandlerFunc {
	return handler.Handle()
}

func middleware(h http.Handler) http.Handler {
	//return logRequest.Handler(h) // MARK: Deactivate forever
	return h
}

func api1(path string) string {
	return fmt.Sprintf("/v1%s", path)
}
