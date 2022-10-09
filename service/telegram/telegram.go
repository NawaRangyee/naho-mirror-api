package telegram

import (
	"fmt"
	"mirror-api/config"
	"mirror-api/util/logger"
	"time"

	tb "gopkg.in/telebot.v3"
)

var to *MonitorRoom
var enabled = false

type MonitorRoom struct{}

func (MonitorRoom) Recipient() string {
	return config.TelegramChatID
}

func getToken() string {
	return config.TelegramAccessToken
}

var bot *tb.Bot

func Init() {
	logger.L.Info("Initializing telegram bot..")
	if config.TelegramChatID == "" {
		logger.L.Fatal("TelegramChatID missing.")
	}
	if config.TelegramAccessToken == "" {
		logger.L.Fatal("TelegramAccessToken missing.")
	}

	var err error

	bot, err = tb.NewBot(tb.Settings{
		Token:  getToken(),
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		logger.L.Fatal(err.Error(), "func", "Init()", "extra", "tb.NewBot()")
	}

	to = &MonitorRoom{}
	enabled = true
}

func SendMessage(message string) {
	SendMessageAt(message, getNow())
}

func SendMessageAt(message string, at time.Time) {
	if !config.IsProductionMode() || !enabled {
		return
	}
	msg := fmt.Sprintf("<%s> %s\n%s", config.ServerName, message, at.In(config.Loc).String())
	logger.L.Debug("Sending telegram Message...")
	_, err := bot.Send(to, msg)
	if err != nil {
		logger.L.Errorw(err.Error(), "func", "SendMessageAt()", "extra", "bot.Send(to, msg)", "to", to.Recipient(), "msg", msg)
		return
	}
	logger.L.Debug("[Telegram] message sent.")
}

func getNow() time.Time {
	return time.Now().In(config.Loc)
}
