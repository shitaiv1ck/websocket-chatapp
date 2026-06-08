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
│  │  └─ transport
│  │     ├─ middleware
│  │     ├─ repsponse
│  │     │  ├─ dto.go
│  │     │  └─ response.go
│  │     └─ request
│  └─ features
├─ Makefile
├─ migrations
│  ├─ 000001_init.down.sql
│  └─ 000001_init.up.sql
└─ readme.md

```