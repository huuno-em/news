# Агрегатор новостей (часть итоговой аттестации)

Для тестирования сервиса отдельно, можно воспользоваться dockerfile, который есть в проекте:
``` go
docker build -t my-postgres .  
docker run -d --name my-postgres-container -p 5432:5432 my-postgres
go run server.go
```
По умолчанию порт для работы сервиса 8083, чтобы изменить его отредактируйте строку 38 в файле server.go
``` go
apiPort := "8083"
```

*Если же запускаем в связке с сервисами [Комменатриев](https://github.com/huuno-em/comments) , [Верификацией](https://github.com/huuno-em/verification) и [APIGateway](https://github.com/huuno-em/api_gateway) , то нужно воспользоваться docker-compose из сервиса комменатриев.*
