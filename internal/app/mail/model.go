package mail

import (
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/emersion/go-imap/client"
)

// Mail - main structure of package
type Mail struct {
	lc           string                 // logging category
	client       *client.Client         // mail client
	imapCfg      config.IMAP            // IMAP config
	smtpCfg      config.SMTP            // SMTP config
	lastMsgID    uint32                 // last id of the read message
	serviceMask  string                 // mask service name
	serviceID    string                 // service id
	chInlineData chan map[string][]byte // channel for receiving messages without attachments
	chCritError  chan bool              // error's handler
}

// subjectData - name of the message description in JSON format
type subjectData struct {
	Description string `json:"description"` // name of the message
	UUID        string `json:"uuid"`        // service id
}
