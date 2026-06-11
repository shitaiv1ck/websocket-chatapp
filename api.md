# API References

**Base URL:** `http://localhost:8080/`

## API ENDPOINTS

| Эндпоинт                              | Описание                              |
| ------------------------------------- | ------------------------------------- |
| /ws                                   | Открывает websocket соединение |
| `POST` /api/users                     | Регистрирует нового пользователя |
| `GET` /api/users?limit=&offset        | Возвращает список пользователей с заданными параметрами `limit` и `offset` (если не заданы, то возвращает всех пользователей) |
| `GET` /api/protected/users/me         | Возвращает текущую сессию пользователя |
| `PATCH` /api/protected/users          | Изменяет логин или пароль пользователя |
| `POST` /api/sessions                  | Создает новую сессию для пользователя |
| `DELETE` /api/protected/sessions      | Удаляет текущую сессию пользователя |
| `POST` /api/protected/friend-requests | Отправляет заявку в друзья пользователю |
| `GET` /api/protected/friend-requests?direction= | Возвращает исходящие (`direction=outgoing`) или входящие (`direction=incoming` или не задано) заявки в друзья |
| `DELETE` /api/protected/friend-requests/{friend_request_id} | Отклоняет заявку в друзья |

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

| Код | Описание | В каких эндпоинтах |
|-----|----------|---------------------|
| 200 | Успешный GET/PATCH | `GET /api/users`, `GET /api/protected/users/me`, `PATCH /api/protected/users`, `GET /api/protected/friend-requests` |
| 201 | Успешное создание | `POST /api/users`, `POST /api/sessions`, `POST /api/protected/friend-requests` |
| 204 | Успешное удаление | `DELETE /api/protected/sessions`, `DELETE /api/protected/friend-requests` |
| 400 | Ошибка валидации | `POST /api/users`, `PATCH /api/protected/users`, `POST /api/sessions`, `POST /api/protected/friend-requests`, `DELETE /api/protected/friend-requests` |
| 401 | Не авторизован | Все `/api/protected/*`, `POST /api/sessions` |
| 404 | Ресурс не найден | `POST /api/protected/friend-requests` (если `to_user_id` не существует), `DELETE /api/protected/friend-requests` (если заявка не найдена) |
| 409 | Конфликт (уже существует) | `POST /api/users`, `PATCH /api/protected/users`, `POST /api/protected/friend-requests` (заявка уже отправлена) |
| 500 | Внутренняя ошибка | Любой |

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
    "username": "some username"
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
        "username": "n1x"
    },
    "to_user": {
        "id": 6,
        "username": "Марк Аврелий"
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
            "username": "n1x"
        },
        "to_user": {
            "id": 1,
            "username": "keynysiro"
        },
        "created_at": "2026-06-11T07:28:08.255482Z"
    },
    {
        "id": 6,
        "from_user": {
            "id": 2,
            "username": "n1x"
        },
        "to_user": {
            "id": 5,
            "username": "shitai"
        },
        "created_at": "2026-06-11T11:51:48.460871Z"
    },
    {
        "id": 7,
        "from_user": {
            "id": 2,
            "username": "n1x"
        },
        "to_user": {
            "id": 6,
            "username": "Марк Аврелий"
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

