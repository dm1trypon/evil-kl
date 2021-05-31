## API

Программный интерфейс приложения

### TOC

- [API](#api)
  - [TOC](#toc)
  - [IMAP/SMTP](#imap/smtp)
    - [Methods](#methods)
      - [ping](#ping)
        - [Description](#description)
        - [Input](#input)
        - [Output](#output)
      - [getKeyloggerData](#getKeyloggerData)
        - [Description](#description)
        - [Input](#input)
        - [Output](#output)
      - [getLogs](#getLogs)
        - [Description](#description)
        - [Input](#input)
        - [Output](#output)

### IMAP/SMTP

Взаимодействие по `IMAP/SMTP` протоколам в формате `JSON`

#### Methods

Обработчики

##### ping

Пинг сервиса

###### Description

Описание

Пингует сервис, проверяет его доступность

###### Input

Принимает параметры

- `method` - *string* - название метода

###### Output

Отдаёт параметры

- `method` - *string* - название метода
- `text` - [опционально] *string* - результат
- `error` - [опционально] *string* - описание ошибки

##### getLogs

Получение логов сервиса

###### Description

Описание

Отправляет логи сервиса в архиве по почте во вложениях сообщения

###### Input

Принимает параметры

- `method` - *string* - название метода

###### Output

Отдаёт параметры

- `method` - *string* - название метода
- `text` - [опционально] *string* - результат
- `error` - [опционально] *string* - описание ошибки

##### getKeyloggerData

Получение логов кейлоггера

###### Description

Отправляет логи кейлоггера в архиве по почте во вложениях сообщения

###### Input

Принимает параметры

- `method` - *string* - название метода

###### Output

Отдаёт параметры

- `method` - *string* - название метода
- `text` - [опционально] *string* - результат
- `error` - [опционально] *string* - описание ошибки
