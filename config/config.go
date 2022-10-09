package config

import (
	"chapi/util/logger"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	ModeProduction  = "Production"
	ModeDevelopment = "Development"

	ServerName = "ch-api"
)

var SqlURL = ""
var MariaDBURL = ""
var TokyoMongoURI = ""
var NodeName = ""
var TelegramChatID = ""
var TelegramAccessToken = ""
var Mode = ""
var UserTokenSecret = ""
var ListenPort = ""
var EmailUsername = ""
var EmailPass = ""
var SimpleSession = ""
var ApiLayerKey = ""
var Loc *time.Location

// MARK: For 3rd-party SDKs

var SDKKakaoClientID = ""
var SDKKakaoClientSecret = ""

func init() {
	envFileLoaded := false
	for i := range os.Args {
		if os.Args[i] == "env" {
			fileName := os.Args[i+1]
			err := godotenv.Load(fileName)
			if err != nil {
				logger.Fatal("Error loading .env file", err)
			}
			envFileLoaded = true
			break
		}
	}

	if !envFileLoaded {
		err := godotenv.Load()
		if err != nil {
			logger.Fatal("Error loading .env file")
		} else {
			logger.Info("Loaded LOT from env file")
		}
	}

	UserTokenSecret = os.Getenv("CHAPI_ENCRYPT")
	if UserTokenSecret == "" {
		logger.Fatal("CHAPI_ENCRYPT missing")
	}

	//MqttURL = os.Getenv("CHAPI_MQTT_URL")
	//if MqttURL == "" {
	//	logger.Fatal("CHAPI_MQTT_URL missing")
	//}
	//MqttClientID = os.Getenv("CHAPI_MQTT_CLIENT_ID")
	//if MqttClientID == "" {
	//	logger.Fatal("CHAPI_MQTT_CLIENT_ID missing")
	//}

	Mode = os.Getenv("CHAPI_MODE")
	if IsProductionMode() {
		logger.Info("Running CHAPI in Production Mode")
	} else {
		Mode = ModeDevelopment
		logger.Info("Running CHAPI in Development Mode")
		log.SetLevel(log.DebugLevel)
	}
	NodeName = os.Getenv("CHAPI_NODE_NAME")

	TokyoMongoURI = os.Getenv("CHAPI_MONGO_URI")
	SqlURL = os.Getenv("CHAPI_MYSQL_URL")

	TelegramChatID = os.Getenv("CHAPI_TELEGRAM_CHAT_ID")
	TelegramAccessToken = os.Getenv("CHAPI_TELEGRAM_ACCESS_TOKEN")

	EmailUsername = os.Getenv("CHAPI_EMAIL_ID")
	EmailPass = os.Getenv("CHAPI_EMAIL_PASSWORD")

	ListenPort = os.Getenv("CHAPI_PORT")

	SimpleSession = os.Getenv("CHAPI_SIMPLE_SESSION")
	if SimpleSession == "" {
		logger.Fatal("CHAPI_SIMPLE_SESSION missing")
	}

	ApiLayerKey = os.Getenv("APILAYER_API_KEY")

	SDKKakaoClientID = os.Getenv("CHAPI_SDK_KAKAO_CLIENT_ID")
	SDKKakaoClientSecret = os.Getenv("CHAPI_SDK_KAKAO_CLIENT_SECRET")

	Loc, _ = time.LoadLocation("Asia/Seoul")
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
