# evil-kl

Сервис кейлоггер.

## TOC
- [evil-kl](#evil-kl)
  - [TOC](#toc)
  - [Description](#description)
  - [Dependencies](#dependencies)
  - [Topics](#topics)

## Description

Описание

Сервис для логирования нажатий клавиш клавиатуры.
Заупускается в два этапа: **установка** и **непосредственный запуск**.
При **установке** создаются необходимые папки для работы сервиса, указываются в конфиге:
 - рабочая папка сервиса
 - папка для логов событий клавиатуры
 - папка для логов сервиса

Создается параметр в ветке реестре для автозагрузки сервиса, название которой указывается в конфиге.
Далее бинарный файл сервиса копируется в рабочую папку сервиса.  
**Непосредственный запуск**: работа сервиса начинается после повторной загрузки ОС.  
Управление сервисом осуществляется методами, описанных в API по протоколам `SMTP/IMAP`.

## Dependencies

Связи и зависимости сервиса

* github.com/dm1trypon/easy-logger
* github.com/davecgh/go-spew/spew
* golang.org/x/sys/windows
* github.com/emersion/go-imap
*	github.com/emersion/go-imap/client
*	github.com/emersion/go-message/mail
*	github.com/scorredoira/email
* github.com/qri-io/jsonschema
* golang.org/x/sys/windows/registry


## Topics

Топики

* [Инсталляция](doc/install/index.md)
* [Конфгурация](doc/configuration/index.md)
* [API](doc/api/index.md)
* [Errors](doc/error/index.md)
* [Требования](doc/requirement/index.md)