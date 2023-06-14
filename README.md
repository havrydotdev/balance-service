[![codecov](https://codecov.io/gh/gavrylenkoIvan/balance-service/branch/master/graph/badge.svg?token=D0E12OKSEI)](https://codecov.io/gh/gavrylenkoIvan/balance-service)
# Test task avitoTech

<!-- ToC start -->
# Content

1. [Task description](#Task-description)
1. [Implementation](#Implementation)
1. [Endpoints](#Endpoints)
1. [Starting](#Starting)
1. [Testing](#Testing)
1. [Examples](#Examples)
<!-- ToC end -->

# Task description

Develop a microservice for working with users' balance (balance, crediting / debiting / transferring funds).
The service must provide an HTTP API and accept/return requests/responses in JSON format.
Additionally, implement methods for converting the balance and obtaining a list of transactions.
Full description in [TASK](TASK.md).
# Implementation

- Following the REST API design.
- Clean architecture and dependency injection
- Working with framework [labstack/echo](https://github.com/labstack/echo).
- Working with Postgres using [sqlx](https://github.com/jmoiron/sqlx) and writing SQL queries.
- App configuration with [viper](https://github.com/spf13/viper) library.
- Launching with Docker.
- Unit/Integration testing using mocks [testify](https://github.com/stretchr/testify), [mock](https://github.com/golang/mock).

**Project structure:**
```
.
├── internal  // business logic
│   ├── handler     
│   ├── service     
│   └── repository  
├── cmd    
├── pkg       // Importable code (logging and utils) 
│   ├── utils     
│   └── logging           
├── schema    // SQL migrations files
├── configs   // App configs
├── models    // Custom types
├── scripts   // Shell scripts
├── docs      // Swagger documentation
```

# Endpoints

- GET /balance/ - получение баланса пользователя
    - Тело запроса:
        - user_id - уникальный идентификатор пользователя.
  - Параметры запроса:
      - currency - валюта баланса.
- GET /transaction/ - получение транзакций пользователя
    - Тело запроса:
        - user_id - уникальный идентификатор пользователя.
    - Параметры запроса:
        - sort - сортировка списка транзакций.
- POST /top-up/ - пополнение баланса пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя,
        - amount - сумма пополнения в RUB.
- POST /debit/ - списание из баланса пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя,
        - amount - сумма списания в RUB.
- POST /transfer/ - перевод средств на баланс другого пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя, с баланса которого списываются средства,
        - to_id - идентификатор пользователя, на баланс которого начисляются средства,
        - amount - сумма перевода в RUB.
# Starting

```
make build
make run
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate-up
```

# Testing

Локальный запуск тестов:
```
make run-test
```

# Examples

Запросы сгенерированы из Postman для cURL.

### 1. GET  /balance для _user_id=1_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/balance' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 1000
}
```

### 2. GET /balance для _user_id=1_ и _currency=USD_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/balance?currency=USD' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 13.542863492536123
}
```

### 3. GET /transaction для _user_id=1_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    },
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    },
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    }
]
```

### 4. GET /transaction для _user_id=1, sort=date_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction?sort=date' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    },
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    }
]
```

### 5. GET /transaction для _user_id=1, sort=amount_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction?sort=amount' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    },
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    }
]
```

### 6. POST /top-up для _user_id=1, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/top-up' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 1000
}
```

### 7. POST /debit для _user_id=1, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/debit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 0
}
```

### 8. POST /transfer для _user_id=1, to_id=2, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "to_id":2,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 2,
    "balance": 1000
}
```
