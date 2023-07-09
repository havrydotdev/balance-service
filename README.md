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
- Unit/Integration testing using mocks [testify](https://github.com/stretchr/testify), [sqlmock](https://github.com/DATA-DOG/go-sqlmock), [gomock](https://github.com/golang/mock/gomock).

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

- GET /balance/{user_id} - get user`s balance
    - URL variables:
        - user_id - unique user`s id.
    - Query params:
      - currency - convert user`s balance to currency (EUR by default).
- GET /transactions/{user_id} - get user`s transactions
    - URL variables:
        - user_id - unique user`s id.
    - Query params:
        - page
        - limit - number of transactions per page 
        - sort - сортировка списка транзакций.
- POST /top-up/{user_id} - replenishment of the user's balance
    - Request body:
        - user_id - unique user`s id,
        - amount - replenishment amount in EUR.
- POST /debit/{user_id} - write-off from the user's balance
    - Request body:
        - user_id - идентификатор пользователя,
        - amount - replenishment amount in EUR.
- POST /transfer/ - transferring funds to the balance of another user
    - Request body:
        - user_id - id of the user from whose balance funds are debited,
        - to_id - id of the user whose balance the funds are credited to,
        - amount - transfer amount in EUR.
# Starting

## Build docker-compose:
```sh
make compose-build
```

## Start container:
```sh
make compose-up
```

# Testing

To run tests, use:
```
make test
```

# Examples

### 1. GET  /balance for _user_id=1_

**Request:**
```
$ curl --location --request GET 'localhost:8080/balance/1' \
--header 'Content-Type: application/json'
```
**Response body:**
```
{
    "user_id": 1,
    "balance": 4.13
}
```

### 2. GET /balance for _user_id=1_ and _currency=USD_

**Request:**
```
$ curl --location --request GET 'localhost:8080/balance/1?currency=UAH' \
--header 'Content-Type: application/json'
```
**Response body:**
```
{
    "user_id": 1,
    "balance": 165.43
}
```

### 3. GET /transactions for _user_id=1_, _page=1_, _limit=1_, _sort=date_

**Request:**
```
$ curl --location --request GET 'localhost:8080/transactions/1?page=1&limit=1&sort=date' \
--header 'Content-Type: application/json'
```
**Response body:**
```
[
   {
        "id": 1,
        "user_id": 1,
        "amount": 30,
        "operation": "",
        "date": "2023-06-14 02:19:40"
   }
]
```

### 4. GET /transactions for _user_id=2_, _page=1_, _limit=1_, _sort=date_

```
$ curl --location --request GET 'localhost:8080/transactions/2?page=1&limit=2&sort=date' \
--header 'Content-Type: application/json'
```
**Response body:**
```
[
    {
        "id": 2,
        "user_id": 2,
        "amount": 101,
        "operation": "",
        "date": "2023-06-14 02:19:40"
    },
    {
        "id": 3,
        "user_id": 2,
        "amount": 32,
        "operation": "",
        "date": "2023-06-14 02:19:40"
    }
]
```

### 5. POST /top-up for _user_id=1, amount=1000_

**Request:**
```
$ curl --location --request POST 'localhost:8080/top-up' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Response body:**
```
{
    "user_id": 1,
    "balance": 1004.13
}
```

### 6. POST /debit for _user_id=1, amount=1000_

**Request:**
```
$ curl --location --request POST 'localhost:8080/debit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Response body:**
```
{
    "user_id": 1,
    "balance": 4.13
}
```
**But if you try to do it again, there will no enough money to perform debit:**
```
{
    "message": "not enough money to perform purchase"
}
```

### 7. POST /transfer for _user_id=1, to_id=2, amount=1000_

**Request:**
```
$ curl --location --request POST 'localhost:8080/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "to_id":2,
    "amount":1
}'
```
**Response body:**
```
{
    "user_id": 2,
    "balance": 33
}
```
