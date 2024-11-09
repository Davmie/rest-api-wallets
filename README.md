# REST API for wallets

CRUD

POST api/v1/wallet
```
{
    walletId: UUID,
    operationType: DEPOSIT or WITHDRAW,
    amount: 1000
}
```

GET api/v1/wallets/{WALLET_UUID}

## Запуск
`docker-compose up -d`

## Тесты
Unit-тесты для repository: `go test ./internal/wallet/repository/postgres`

Unit-тесты для usecase: `go test ./internal/wallet/usecase/`

Для проверки работоспособности сервера можно импортировать коллекцию в Postman: `./postman/postman_collection.json`