# API ENDPOINTS

| Endpoint                              | Описание                              |
| ------------------------------------- | ------------------------------------- |
| /ws                                   | Открытие websocket соединения         |
| `POST` /api/users                     | Создание нового пользователя          |
| `GET` /api/users                      | Получение всех пользователей          |
| `GET` /api/protected/users/me         | Получение текущей сессии (необходимо пройти аутентификацию) |
| `PATCH` /api/protected/users          | Изменение данных о пользователе       |
| `POST` /api/sessions                  | Создание новой сессии                 |
| `DELETE` /api/protected/sessions      | Удаление текущей сессии               |

Для паттернов группы `/api/protected/` c **небезопасными методами** дополнительно проверятеся аутентификация + csrf токен, переданный в запросе под заголовком **X-CSRF-Token**

## Примеры JSON в теле запроса/ответа

**`POST` /api/users:**

Request body:
```JSON
{
    "username": "some username",  
    "password": "some password"  
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
{}
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
{}
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
    "error": "invalid coockie's files",
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
    "username": "NewUsername",
    "old_password": "some password",
    "new_password": "newpassword"
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
    "error": "failed to create user: already exists",
    "message": "failed to create user"
}
```
```JSON
401 Unauthorized

{
    "error": "invalid coockie's files",
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

{}
```
```JSON
401 Unauthorized

{
    "error": "invalid username or password",
    "message": "failed to authentication"
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
{}
```

Response body:
```JSON
204 No Content

{}
```
```JSON
401 Unauthorized

{
    "error": "invalid coockie's files",
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



