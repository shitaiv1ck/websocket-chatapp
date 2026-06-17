# API References

**Base URL:** `http://localhost:8080/`

## API ENDPOINTS

| Эндпоинт                              | Описание                              |
| ------------------------------------- | ------------------------------------- |
| /ws                                   | Открывает websocket соединение |
| `POST` /api/users                     | Регистрирует нового пользователя |
| `GET` /api/users?search=&limit=&offset=       | Возвращает список пользователей с заданными параметрами `search`, `limit` и `offset` (если не заданы, то возвращает всех пользователей) |
| `GET` /api/users/me         | Возвращает текущую сессию пользователя |
| `PATCH` /api/users          | Изменяет логин или пароль пользователя |
| `POST` /api/sessions                  | Создает новую сессию для пользователя |
| `DELETE` /api/sessions      | Удаляет текущую сессию пользователя |
| `POST` /api/friend-requests | Отправляет заявку в друзья пользователю |
| `GET` /api/friend-requests?direction= | Возвращает исходящие (`direction=outgoing`) или входящие (`direction=incoming` или не задано) заявки в друзья |
| `DELETE` /api/friend-requests/{friend_request_id} | Отклоняет заявку в друзья |
| `POST` /api/friendships | Принимает заявку в друзья |
| `GET` /api/friendships?limit=&offset= | Возвращает список друзей с заданными параметрами `limit` и `offset` (если не заданы, то возвращает всех друзей)|
| `DELETE` /api/friendships/{friendship_id} | Удаляет друга |
| `POST` /api/chats | Создает или возвращает чат с другом |
| `GET` /api/chats?limit=&offset= | Возващает список чатов с заданными параметрами `limit` и `offset` (если не заданы, то возвращает все чаты) |
| `DELETE` /api/chats/{chat_id} | Удаляет чат |
| `POST` /api/chats/{chat_id}/messages | Отправляет сообщение в чат |
| `GET` /api/chats/{chat_id}/messages | Возвращает отсортированный список всех сообщений в чате |


### Примечание 

Решение сделать паттерн `POST /api/chats` идемпотентным было принято, чтобы избежать конфликта, когда пользователь А и пользователь Б одновременно захотели начать новый чат друг с другом

### CSRF Protection

Для методов `POST`, `PATCH`, `DELETE` обязателен заголовок `X-CSRF-Token`, **кроме** паттернов `POST` /api/sessions и  `POST` /api/users

### Аутентификация

После успешного `POST /api/sessions` сервер устанавливает cookie:
- **Name:** `session_token`
- **HttpOnly:** true
- **SameSite:** Lax
- **Max-Age:** 24 часа

И

- **Name:** `csrf_token`
- **HttpOnly:** false
- **SameSite:** Lax
- **Max-Age:** 24 часа

## Коды ответов

| Код | Описание | Эндпоинты |
|-----|----------|-----------|
| 200 | Успешный GET / POST (идемпотентный) | `GET /api/users`, `GET /api/users/me`, `GET /api/friend-requests`, `GET /api/friendships`, `GET /api/chats`, `GET /api/chats/{chat_id}/messages`, `POST /api/chats` (чат уже существовал) |
| 201 | Успешное создание | `POST /api/users`, `POST /api/sessions`, `POST /api/friend-requests`, `POST /api/friendships`, `POST /api/chats/{chat_id}/messages` |
| 204 | Успешное удаление | `DELETE /api/sessions`, `DELETE /api/friend-requests/{id}`, `DELETE /api/friendships/{id}`, `DELETE /api/chats/{id}` |
| 400 | Ошибка валидации | `POST /api/users`, `PATCH /api/users`, `POST /api/sessions`, `POST /api/friend-requests`, `POST /api/friendships`, `POST /api/chats` (попытка создать чат с собой), `POST /api/chats/{chat_id}/messages` (пользователь не друг), `DELETE /api/chats/{id}` (неверный формат ID) |
| 401 | Не авторизован | Все эндпоинты кроме `POST /api/users`, `GET /api/users`, `POST /api/sessions`. Также `GET /api/users/me`, `PATCH /api/users`, `DELETE /api/sessions` и все эндпоинты `/api/friend-requests`, `/api/friendships`, `/api/chats`, `/api/chats/{chat_id}/messages` |
| 404 | Ресурс не найден | `POST /api/friend-requests` (пользователь не существует), `DELETE /api/friend-requests/{id}` (заявка не найдена), `POST /api/friendships` (запрос не найден), `DELETE /api/friendships/{id}` (дружба не найдена), `POST /api/chats` (друг не найден), `DELETE /api/chats/{id}` (чат не найден), `POST /api/chats/{chat_id}/messages` (чат не найден), `GET /api/chats/{chat_id}/messages` (чат не найден) |
| 409 | Конфликт (уже существует) | `POST /api/users` (username занят), `PATCH /api/users` (username занят), `POST /api/friend-requests` (заявка уже отправлена или уже друзья) |
| 500 | Внутренняя ошибка сервера | Любой |

## Примеры JSON в теле запроса/ответа

**`POST` /api/users:**

Request body:
```JSON
{
    "username": "some username", #required, min=3, max=100
    "password": "some password"  #required, min=8, max=100
}
```

Response body:
```JSON
201 Created

{
    "id": 3,
    "username": "some username"
}
```
```JSON
409 Conflict

{
    "error": "failed to create user: already exists",
    "message": "failed to create user"
}
```
```JSON
400 Bad Request

{
    "error": "failed to validate json: invalid argument",
    "message": "failed to decode and validate"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`GET` /api/users:**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

[
    {
        "id": 1,
        "username": "Keyny",
        "is_online": false
    },
    {
        "id": 2,
        "username": "KeynySiro",
        "is_online": false
    },
    {
        "id": 3,
        "username": "some username",
        "is_online": false
    }
]
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`GET` /api/users/me**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

{
    "id": 3,
    "username": "some username",
    "is_online": true
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`PATCH` /api/users**

Request body:
```JSON
{
    "username": "NewUsername",        #optional, not null
    "old_password": "some password",  #optional, not null
    "new_password": "newpassword"     #optional, not null
}
```

Response body:
```JSON
200 OK

{
    "id": 3,
    "username": "NewUsername"
}
```
```JSON
400 Bad Request

{
    "error": "invalid password: invalid argument",
    "message": "failed to patch user"
}
```
```JSON
409 Conflict

{
    "error": "failed to patch user: already exists",
    "message": "failed to patch user"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`POST` /api/sessions**
Request body:
```JSON
{
    "username": "NewUsername",
    "password": "newpassword"
}
```

Response body:
```JSON
201 Created

(no body)
```
```JSON
401 Unauthorized

{
    "error": "invalid username or password",
    "message": "failed to authenticate"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`DELETE` /api/sessions**

Request body:
```JSON
(no body)
```

Response body:
```JSON
204 No Content

(no body)
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`POST` /api/friend-requests**

Request body:
```JSON
{
    "to_user_id": 6  #required
}
```

Response body:
```JSON
201 Created

{
    "id": 7,
    "from_user": {
        "id": 2,
        "username": "n1x",
        "is_online": false
    },
    "to_user": {
        "id": 6,
        "username": "Марк Аврелий",
        "is_online": false
    },
    "created_at": "2026-06-11T15:19:02.616126Z"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```
```JSON
409 Conflict

{
    "error": "already exists",
    "message": "failed to create friend request"
}
```

**`GET` /api/friend-requests**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

[
    {
        "id": 1,
        "from_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "to_user": {
            "id": 1,
            "username": "keynysiro",
            "is_online": false
        },
        "created_at": "2026-06-11T07:28:08.255482Z"
    },
    {
        "id": 6,
        "from_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "to_user": {
            "id": 5,
            "username": "shitai",
            "is_online": false
        },
        "created_at": "2026-06-11T11:51:48.460871Z"
    },
    {
        "id": 7,
        "from_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "to_user": {
            "id": 6,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "created_at": "2026-06-11T15:19:02.616126Z"
    }
]
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`DELETE` /api/friend-requests/{friend_request_id}**

Request body:
```JSON
(no body)
```

Response body:
```JSON
204 No Content

(no body)
```
```JSON
404 Not Found

{
    "error": "failed to delete friend request: friend request doesn't exist: not found",
    "message": "failed to delete friend request"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`POST` /api/friendships**

Request body:
```JSON
{
    "friend_request_id": 8  #required
}
```

Response body:
```JSON
201 Created

{
    "id": 5,
    "first_user": {
        "id": 3,
        "username": "n1x",
        "is_online": false
    },
    "second_user": {
        "id": 5,
        "username": "shitai",
        "is_online": false
    }
}
```
```JSON
404 Not Found

{
    "error": "failed to get friend request from rep: not found",
    "message": "failed to create frienship"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`GET` /api/friendships**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

[
    {
        "id": 2,
        "first_user": {
            "id": 1,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "second_user": {
            "id": 3,
            "username": "n1x",
            "is_online": false
        }
    },
    {
        "id": 4,
        "first_user": {
            "id": 3,
            "username": "n1x",
            "is_online": false
        },
        "second_user": {
            "id": 4,
            "username": "KeynySiro",
            "is_online": false
        }
    }
]
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`DELETE` /api/friendships/{friendships_id}**

Request body:
```JSON
(no body)
```

Response body:
```JSON
204 No Content

(no body)
```
```JSON
404 Not Found

{
    "error": "failed to delete friendship: friendship doesn't exist: not found",
    "message": "failed to delete friendship"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`POST` /api/chats**

Request body:
```JSON
{
    "friend_id": 3  #required
}
```

Response body:
```JSON
200 OK

{
    "id": 7,
    "first_user": {
        "id": 2,
        "username": "n1x",
        "is_online": false
    },
    "second_user": {
        "id": 3,
        "username": "Марк Аврелий",
        "is_online": false
    },
    "last_message_content": null,
    "last_message_at": "2026-06-14T15:37:57.62259+03:00"
}
```
```JSON
400 Bad Request

{
    "error": "can't create chat with yourself: invalid argument",
    "message": "failed to create or get chat"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`GET` /api/chats**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

[
    {
        "id": 7,
        "first_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "second_user": {
            "id": 3,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "last_message_content": null,
        "last_message_at": "2026-06-14T15:37:57.62259+03:00"
    },
    {
        "id": 4,
        "first_user": {
            "id": 1,
            "username": "KeynySiro",
            "is_online": false
        },
        "second_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "last_message_content": null,
        "last_message_at": "2026-06-14T15:37:47.239136+03:00"
    }
]
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`DELETE` /api/chats/{chat_id}**

Request body:
```JSON
(no body)
```

Response body:
```JSON
204 No Content

(no body)
```
```JSON
404 Not Found

{
    "error": "failed to delete chat: chat doesn't exist: not found",
    "message": "failed to delete chat"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`POST` /api/chats/{chat_id}/messages**

Request body:
```JSON
{
    "receiver_id": 3,  #required
    "content": "Hi!"   #required
}
```

Response body:
```JSON
201 Created

{
    "id": 11,
    "chat_id": 1,
    "sender_id": 1,
    "receiver_id": 3,
    "content": "Hi!",
    "created_at": "2026-06-15T09:39:01.869267+03:00"
}
```
```JSON
400 Bad Request

{
    "error": "user with id=2 isn't your friend: invalid argument",
    "message": "failed to create message"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

**`GET` /api/chats/{chat_id}/messages**

Request body:
```JSON
(no body)
```

Response body:
```JSON
200 OK

[
    {
        "id": 7,
        "chat_id": 1,
        "sender_id": 3,
        "receiver_id": 1,
        "content": "Hi brooo!",
        "created_at": "2026-06-14T22:07:00.826594+03:00"
    },
    {
        "id": 8,
        "chat_id": 1,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Oh, my tarnished son...",
        "created_at": "2026-06-14T22:09:05.370716+03:00"
    }
]
```
```JSON
401 Unauthorized

{
    "error": "invalid cookie",
    "message": "check cookie"
}
```
```JSON
500 Internal Server Error

{
    "error": "some error",
    "message": "some message"
}
```

## WebSocket соединение

### Проверка соединения

| Параметр | Значение | Описание |
|----------|----------|----------|
| `pingPeriod` | 54 секунды | Интервал отправки ping фреймов |
| `pongWait` | 60 секунд | Максимальное время ожидания pong |

Сервер поддерживает соединение через автоматический ping/pong механизм:

1. Сервер отправляет **протокольный ping фрейм** каждые 54 секунды
2. Клиент должен ответить **pong фреймом** в течение 60 секунд
3. При превышении таймаута соединение закрывается.

### Сообщения от сервера (примеры)

- Новый зарегистрированный пользователь (получают все пользователи):

```JSON
{
    "type": "user.created",
    "content": {
        "id": 1,
        "username": "KeynySiro",
        "is_online": false
    }
}
```

- Пользователь стал онлайн/оффлайн (получают все пользователи):

```JSON
{
    "type": "user.change_status",
    "content": {
        "id": 1,
        "username": "KeynySiro",
        "is_online": true
    }
}
```

- Пользователь изменил username (получают все пользователи):

```JSON
{
    "type": "user.change_username",
    "content": {
        "id": 1,
        "username": "NewUsername",
        "is_online": true
    }
}
```

- Пользователь отправил заявку в друзья (получает тот, кому отправили заявку в друзья):

```JSON
{
    "type": "friend_request.received",
    "content": {
        "friend_request_id": 2,
        "from_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "to_user": {
            "id": 6,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "created_at": "2026-06-11T15:19:02.616126Z"
    }
}
```

- Пользователь отправил заявку в друзья (получает тот, кто отправил заявку в друзья):

```JSON
{
    "type": "friend_request.sent",
    "content": {
        "friend_request_id": 2,
        "from_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "to_user": {
            "id": 6,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "created_at": "2026-06-11T15:19:02.616126Z"
    }
}
```

- Пользователь отклонил заявку в друзья (получают оба пользователя):

```JSON
{
    "type": "friend_request.declined",
    "content": {
        "friend_request_id": 2
    }
}
```

- Пользователь принял заявку в друзья (получает тот, кто принял заявку):

```JSON
{
    "type": "friend_request.accepted",
    "content": {
        "friend_request_id":2,
        "friendship_id": 5,
        "first_user": {
            "id": 3,
            "username": "n1x",
            "is_online": false
        },
        "second_user": {
            "id": 5,
            "username": "shitai",
            "is_online": false
        }
    }
}
```

- Пользователь принял заявку в друзья (получает тот, чью заявку приняли):

```JSON
{
    "type": "friendship.added",
    "content": {
        "friend_request_id":2,
        "friendship_id": 5,
        "first_user": {
            "id": 3,
            "username": "n1x",
            "is_online": false
        },
        "second_user": {
            "id": 5,
            "username": "shitai",
            "is_online": false
        }
    }
}
```

- Пользователь удалил друга (получают оба пользователя):

```JSON
{
    "type": "friendship.deleted",
    "content": {
        "friendship_id": 2
    }
}
```

- Пользователь создал чат (получают оба пользователя):

```JSON
{
    "type": "chat.created",
    "content": {
        "id": 7,
        "first_user": {
            "id": 2,
            "username": "n1x",
            "is_online": false
        },
        "second_user": {
            "id": 3,
            "username": "Марк Аврелий",
            "is_online": false
        },
        "last_message_content": null,
        "last_message_at": "2026-06-14T15:37:57.62259+03:00"
    }
}
```

- Пользователь удалил чат (получают оба пользователя):

```JSON
{
    "type": "chat.deleted",
    "content": {
        "chat_id": 2
    }
}
```

- Пользователь отправил сообщение в чат (получает тот, кто отправил сообщение):

```JSON
{
    "type": "message.sent",
    "content": {
        "id": 11,
        "chat_id": 1,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hi!",
        "created_at": "2026-06-15T09:39:01.869267+03:00"
    }
}
```

- Пользователь отправил сообщение в чат (получает тот, кому отправили сообщение):

```JSON
{
    "type": "message.received",
    "content": {
        "id": 11,
        "chat_id": 1,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hi!",
        "created_at": "2026-06-15T09:39:01.869267+03:00"
    }
}
```