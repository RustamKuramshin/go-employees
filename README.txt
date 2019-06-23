Go Employee.

REST-сервис на Golang, который позволяет редактировать список сотрудников. В качестве СУБД используется MySQL.

Описание endpoints:
1) Создать сотрудника:
> curl -X POST \
  'http://localhost:8080/employee?name=Joe' \
  -H 'Accept: */*' \
  -H 'Authorization: Basic YWRtaW46YWRtaW4=' \
  -H 'Host: localhost:8080' \
  -H 'accept-encoding: gzip, deflate'

2) Получить всех сотрудников:
> curl -X GET \
  http://localhost:8080/employee \
  -H 'Accept: */*' \
  -H 'Authorization: Basic YWRtaW46YWRtaW4=' \
  -H 'Host: localhost:8080' \
  -H 'accept-encoding: gzip, deflate'

3) Изменить имя сотрудника по id:
> curl -X PUT \
  'http://localhost:8080/employee/1?name=John' \
  -H 'Accept: */*' \
  -H 'Authorization: Basic YWRtaW46YWRtaW4=' \
  -H 'Host: localhost:8080' \
  -H 'accept-encoding: gzip, deflate'
4) Удалить сотрудника по имени:
> curl -X DELETE \
  'http://localhost:8080/employee?name=John' \
  -H 'Accept: */*' \
  -H 'Authorization: Basic YWRtaW46YWRtaW4=' \
  -H 'Host: localhost:8080' \
  -H 'accept-encoding: gzip, deflate'

Запуск: 
> docker-compose up --build
