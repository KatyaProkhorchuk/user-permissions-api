# user-permissions-api


### Проблема:
В нашем продукте есть несколько микросервисов. Для работы каждого из них требуется знать права пользователя, который запустил ту или иную задачу. На архитектурном комитете приняли решение централизовать работу с привилегиями пользователя и вынести в отдельный сервис.

### Задача:
Необходимо реализовать микросервис для работы с привилегиями пользователей (создание/ удаление пользователя, добавить/убрать права пользователя, проверка прав пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.

Для запуска сервиса 

```
cd user-permissions-api
docker-compose up --build access
```

Для запуска тестов

```
cd user-permissions-api
docker-compose up --build test
```

Выполним несколько запросов(все тестовые файлы находятся в `exampleRequests`)

### Добавление нового пользовтеля

```
curl http://localhost:4321/insertUser -H 'Content-Type:application/json' -d @exampleRequests/insert.json
```
Ответ сервера

```
"User successfuly created"
```

### Добавление прав пользователя

```
curl http://localhost:4321/addUserRights -H 'Content-Type:application/json' -d @exampleRequests/addRights.json
```

### Проверим права пользователя

```
curl -d @exampleRequests/checkAccess.json -H 'Content-Type:application/json' http://localhost:1234/checkAccess
```

Ответ сервера

```
{"access":["agent2","newagent1","newagent2"]}
```
PS В `checkAccess.json` нужен токен, что бы его получить выполним 

```
curl -X POST   http://localhost:6789/get_token   -H 'Content-Type: application/x-www-form-urlencoded'   -d 'grant_type=client_credentials&client_id=111111&client_secret=12345'

```

Ответ сервера

```
{"access_token":"ZJM0NMZJZWITNMY1OC0ZMJE1LWFJNJUTYZDJMMUXYMU4ODC0","expires_in":7200,"token_type":"Bearer"}
```
### Удаление прав пользователя

```
curl http://localhost:4321/deleteUserRights -H 'Content-Type:application/json' -d @exampleRequests/deleteRights.json
```

Ответ сервера с правами доступа

```
{"access":["agent2","newagent1"]}
```

