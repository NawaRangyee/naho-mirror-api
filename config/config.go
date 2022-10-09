package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"mirror-api/util/logger"
	"os"
	"strings"
	"time"
)

const (
	ModeProduction  = "Production"
	ModeDevelopment = "Development"

	ServerName = "mirror-api"
)

var NodeName = ""
var TelegramChatID = ""
var TelegramAccessToken = ""
var Mode = ""
var ListenPort = ""
var Loc *time.Location

var LogPath = flag.String("log-path", "./logs", "log path for saving locally")
var dotEnvPath = flag.String("env", ".env", "env file to specify. (default: '.env')")

func init() {
	err := godotenv.Load(*dotEnvPath)
	if err != nil {
		logger.Fatalln("Error loading .env file")
	} else {
		logger.Infoln("Loaded LOT from env file")
	}

	//MqttURL = os.Getenv("CHAPI_MQTT_URL")
	//if MqttURL == "" {
	//	logger.Fatalln("CHAPI_MQTT_URL missing")
	//}
	//MqttClientID = os.Getenv("CHAPI_MQTT_CLIENT_ID")
	//if MqttClientID == "" {
	//	logger.Fatalln("CHAPI_MQTT_CLIENT_ID missing")
	//}

	Mode = os.Getenv("MODE")
	if IsProductionMode() {
		logger.Infoln("Running mirror-api in Production Mode")
	} else {
		Mode = ModeDevelopment
		logger.Infoln("Running mirror-api in Development Mode")
		log.SetLevel(log.DebugLevel)
	}
	NodeName = os.Getenv("NODE_NAME")

	TelegramChatID = os.Getenv("TELEGRAM_CHAT_ID")
	TelegramAccessToken = os.Getenv("TELEGRAM_ACCESS_TOKEN")

	ListenPort = os.Getenv("LISTEN_PORT")

	Loc, _ = time.LoadLocation("Asia/Seoul") // Change this value to your location
}

func IsProductionMode(strs ...string) bool {
	if len(strs) > 0 {
		return strings.EqualFold(strs[0], ModeProduction)
	}
	return strings.EqualFold(Mode, ModeProduction)
}

func GetServiceFullName() string {
	return fmt.Sprintf("%s-%s", ServerName, NodeName)
}
