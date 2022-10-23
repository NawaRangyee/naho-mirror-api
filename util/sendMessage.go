package util

import (
	"fmt"
	"mirror-api/config"
	"mirror-api/service/telegram"
	"mirror-api/util/logger"
	"time"
)

func SendFailed(location string, err error) {
	at := time.Now().In(config.Loc)
	msg := fmt.Sprintf("[ERROR/%s]\n=> %s", location, err)
	logger.L.Errorw(err.Error(), "location", location)

	telegram.SendMessageAt(msg, at)
}

func SendNotice(header, location, content string) {
	at := time.Now().In(config.Loc)
	msg := fmt.Sprintf("[%s/%s]\n=> %s\n", header, location, content)
	logger.L.Infow(content, "header", header, "location", location)

	telegram.SendMessageAt(msg, at)
}

func SendStarted(hostname string, localIP string, pubIP string) {
	msg := fmt.Sprintf("Server started successfully\nHostname:%s\nLocal IP:%s\nPublic IP:%s", hostname, localIP, pubIP)
	logger.L.Infow(msg, "func", "SendStarted()", "hostname", hostname, "localIP", localIP, "pubIP", pubIP)

	telegram.SendMessage(msg)
}

func SendStopped(hostname string, localIP, pubIP string) {
	msg := fmt.Sprintf("Server stopping normally\nHostname:%s\nLocal IP:%s\nPublic IP:%s", hostname, localIP, pubIP)
	logger.L.Infow(msg, "func", "SendStopped()", "hostname", hostname, "localIP", localIP, "pubIP", pubIP)

	telegram.SendMessage(msg)
}

func SendForcedStopped(hostname string, localIP, pubIP string) {
	msg := fmt.Sprintf("Server stopping forcingly\nHostname:%s\nLocal IP:%s\nPublic IP:%s", hostname, localIP, pubIP)
	logger.L.Warnw(msg, "func", "SendForcedStopped()", "hostname", hostname, "localIP", localIP, "pubIP", pubIP)

	telegram.SendMessage(msg)
}
