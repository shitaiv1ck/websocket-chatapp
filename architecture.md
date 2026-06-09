# Структура проекта

```
chat
├─ architecture.md
├─ cmd
│  └─ chat
│     └─ main.go
├─ docker-compose.yaml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ core
│  │  ├─ domains
│  │  │  ├─ nullable.go
│  │  │  └─ user.go
│  │  ├─ errors
│  │  │  └─ errors.go
│  │  ├─ logger
│  │  │  ├─ config.go
│  │  │  └─ logger.go
│  │  ├─ server
│  │  │  ├─ http
│  │  │  │  ├─ config.go
│  │  │  │  └─ server.go
│  │  │  └─ ws
│  │  │     ├─ client.go
│  │  │     ├─ message.go
│  │  │     └─ server.go
│  │  ├─ store
│  │  │  └─ postgres
│  │  │     ├─ config.go
│  │  │     └─ postgres.go
│  │  ├─ transport
│  │  │  ├─ middleware
│  │  │  │  ├─ common.go
│  │  │  │  ├─ middleware.go
│  │  │  │  └─ protected.go
│  │  │  ├─ repsponse
│  │  │  │  ├─ dto.go
│  │  │  │  └─ response.go
│  │  │  └─ request
│  │  │     ├─ decode.go
│  │  │     └─ pathvalue.go
│  │  └─ utils
│  │     └─ context.go
│  └─ features
│     └─ users
│        ├─ repository
│        │  └─ repository.go
│        ├─ service
│        │  └─ service.go
│        └─ transport
│           ├─ http
│           │  ├─ dto.go
│           │  └─ transport.go
│           └─ ws
│              ├─ dto.go
│              └─ transport.go
├─ Makefile
├─ migrations
│  ├─ 000001_init.down.sql
│  └─ 000001_init.up.sql
└─ readme.md

```

# Схема БД (5 таблиц)

| Таблица            | Первичный ключ                 | Описание                                                    |
| -----------------  | ------------------------------ | ----------------------------------------------------------- |
| **users**          | `id` (GENERATED)               | Пользователь (хранит информацию о пользователе)             |
| **sessions**       | `session_token` (VARCHAR(255)) | Сессия пользователя (хранит файлы cookie + ID пользователя) |
| **friendships**    | `id` (GENERATED)               | Друзья (хранит информацию о дружественных связях между пользователями) |
| **friendrequests** | `id` (GENERATED)               | Заявки в друзья (хранит информацию о том, кто и кому отправил звявку в друзья) |
| **messages**       | `id` (GENERATED)               | Сообщениея (хранит информацию об отправленных сообщениях) |

## Диаграмма

![diagarm](./docs/diagram.png)