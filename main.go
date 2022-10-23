package main

import (
	"fmt"
	"github.com/chyeh/pubip"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"mirror-api/config"
	"mirror-api/controller/route"
	"mirror-api/data"
	"mirror-api/docs"
	"mirror-api/model/kvDB/mirror"
	"mirror-api/service/checkRsync"
	"mirror-api/service/telegram"
	"mirror-api/util"
	"mirror-api/util/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var localIP = ""
var pubIP = ""
var hostname = ""

var c *cron.Cron

func init() {
	logger.Init(config.IsProductionMode())
	data.Init()

	// MARK: Init model
	mirror.InitFromConfig()

	// MARK: Init Services...
	telegram.Init()
}

func main() {
	// MARK: Send Startup Message
	go sendStarted()

	// MARK: Loading HTTP Router
	router := route.Load()

	// MARK: Allowing nCors for now for frontend reset pass
	corsMiddle()
	nCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		MaxAge:           1728000,
	})
	handler := nCors.Handler(router)

	// MARK: programatically set swagger info
	docs.SwaggerInfo.Title = "Mirror API"
	docs.SwaggerInfo.Description = "Access-Token : AccessToken"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.GetHostname()
	docs.SwaggerInfo.BasePath = "/v1"

	// MARK: Starting cron jobs
	cronJobs()
	// MARK: Handling server requested to stop.
	handleServerStop()

	// MARK: Run once before API Server starts
	go checkRsync.Run()

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

func cronJobs() {
	if !config.IsProductionMode() {
		return
	}
	if c == nil {
		c = cron.New(cron.WithLocation(time.Local))
	}

	// MARK: Renew config.json every 10 minutes
	_, _ = c.AddFunc("*/10 * * * *", mirror.InitFromConfig)
	_, _ = c.AddFunc("* * * * *", checkRsync.Run)

	c.Start()
}
