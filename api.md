# API References

**Base URL:** `http://localhost:8080/`

## API ENDPOINTS

| Эндпоинт                              | Описание                              |
| ------------------------------------- | ------------------------------------- |
| /ws                                   | Открывает websocket соединение |
| `POST` /api/users                     | Регистрирует нового пользователя |
| `GET` /api/users?limit=&offset=       | Возвращает список пользователей с заданными параметрами `limit` и `offset` (если не заданы, то возвращает всех пользователей) |
| `GET` /api/protected/users/me         | Возвращает текущую сессию пользователя |
| `PATCH` /api/protected/users          | Изменяет логин или пароль пользователя |
| `POST` /api/sessions                  | Создает новую сессию для пользователя |
| `DELETE` /api/protected/sessions      | Удаляет текущую сессию пользователя |
| `POST` /api/protected/friend-requests | Отправляет заявку в друзья пользователю |
| `GET` /api/protected/friend-requests?direction= | Возвращает исходящие (`direction=outgoing`) или входящие (`direction=incoming` или не задано) заявки в друзья |
| `DELETE` /api/protected/friend-requests/{friend_request_id} | Отклоняет заявку в друзья |
| `POST` /api/protected/friendships | Принимает заявку в друзья |
| `GET` /api/protected/friendships?limit=&offset= | Возвращает список друзей с заданными параметрами `limit` и `offset` (если не заданы, то возвращает всех друзей)|
| `DELETE` /api/protected/friendships/{friendship_id} | Удаляет друга |
| `POST` /api/protected/chats | Создает или возвращает чат с другом |
| `GET` /api/protected/chats?limit=&offset= | Возващает список чатов с заданными параметрами `limit` и `offset` (если не заданы, то возвращает все чаты) |
| `DELETE` /api/protected/chats/{chat_id} | Удаляет чат |

### Примечание 

Решение сделать паттерн `POST /api/protected/chats` идемпотентным было принято, чтобы избежать конфликта, когда пользователь А и пользователь Б одновременно захотели начать новый чат друг с другом

### CSRF Protection

 Для паттернов группы `/api/protected/*` обязательна пройденная аутентификация. Помимо этого, для методов `POST`, `PATCH`, `DELETE` в этой же группе обязателен заголовок `X-CSRF-Token`

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
| 200 | Успешный GET / POST (идемпотентный) | `GET /api/users`, `GET /api/protected/users/me`, `GET /api/protected/friend-requests`, `GET /api/protected/friendships`, `GET /api/protected/chats`, `POST /api/protected/chats` (чат уже существовал) |
| 201 | Успешное создание | `POST /api/users`, `POST /api/sessions`, `POST /api/protected/friend-requests`, `POST /api/protected/friendships` |
| 204 | Успешное удаление | `DELETE /api/protected/sessions`, `DELETE /api/protected/friend-requests/{id}`, `DELETE /api/protected/friendships/{id}`, `DELETE /api/protected/chats/{id}` |
| 400 | Ошибка валидации | `POST /api/users`, `PATCH /api/protected/users`, `POST /api/sessions`, `POST /api/protected/friend-requests`, `POST /api/protected/friendships`, `POST /api/protected/chats` (попытка создать чат с собой), `DELETE /api/protected/chats/{id}` (неверный формат ID) |
| 401 | Не авторизован | Все `/api/protected/*`, `POST /api/sessions` |
| 404 | Ресурс не найден | `POST /api/protected/friend-requests` (пользователь не существует), `DELETE /api/protected/friend-requests/{id}` (заявка не найдена), `POST /api/protected/friendships` (запрос не найден), `DELETE /api/protected/friendships/{id}` (дружба не найдена), `POST /api/protected/chats` (друг не найден), `DELETE /api/protected/chats/{id}` (чат не найден) |
| 409 | Конфликт (уже существует) | `POST /api/users` (username занят), `PATCH /api/protected/users` (username занят), `POST /api/protected/friend-requests` (заявка уже отправлена или уже друзья) |
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

**`GET` /api/protected/users/me**

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

**`PATCH` /api/protected/users**
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

`DELETE` /api/protected/sessions

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

**`POST` /api/protected/friend-requests**

Request body:
```JSON
{
    "to_user_id": 6
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

**`GET` /api/protected/friend-requests**

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

**`DELETE` /api/protected/friend-requests/{friend_request_id}**

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

**`POST` /api/protected/friendships**

Request body:
```JSON
{
    "friend_request_id": 8
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

**`GET` /api/protected/friendships**

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

**`DELETE` /api/protected/friendships/{friendships_id}**

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

**`POST` /api/protected/chats**

Request body:
```JSON
{
    "friend_id": 3
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

**`GET` /api/protected/chats**

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

**`DELETE` /api/protected/chats/{chat_id}**

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