package main

import (
	"fmt"
	"github.com/chyeh/pubip"
	"github.com/rs/cors"
	"mirror-api/config"
	"mirror-api/docs"
	"mirror-api/service/telegram"
	"mirror-api/util"
	"mirror-api/util/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var localIP = ""
var pubIP = ""
var hostname = ""

func init() {
	logger.Init(config.IsProductionMode())
	data.Init()

	// Init Services...
	telegram.Init()
}

func main() {
	// Send Startup Message
	go sendStarted()

	// Loading HTTP Router
	router := route.Load()

	//Allowing c for now for frontend reset pass
	corsMiddle()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		MaxAge:           1728000,
	})
	handler := c.Handler(router)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Mirror API"
	docs.SwaggerInfo.Description = "Access-Token : AccessToken"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.GetHostname()
	docs.SwaggerInfo.BasePath = "/v1"

	// Handling server requested to stop.
	handleServerStop()

	logger.L.Infof("Listening on port %s\n", config.ListenPort)
	logger.L.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ListenPort), handler))
}

func corsMiddle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(nil)
		}
	})
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.L.Fatalw(err.Error(), "func", "GetOutboundIP()", "extra", `net.Dial("udp", "8.8.8.8:80")`)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func handleServerStop() {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			util.SendStopped(hostname, localIP, pubIP)
			os.Exit(0)
		}
		if sig == syscall.SIGKILL {
			util.SendForcedStopped(hostname, localIP, pubIP)
			os.Exit(0)
		}
	}()
}

func sendStarted() {
	hostname, _ = os.Hostname()
	localIP = GetOutboundIP().String()
	pIP, _ := pubip.Get()
	pubIP = pIP.String()
	if pIP == nil {
		pubIP = localIP
		localIP = "<nil>"
	}
	// Send a telegramMessage to notice server has been started.
	util.SendStarted(hostname, localIP, pubIP)
}
