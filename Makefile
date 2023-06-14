build:
	go build -o balance ./cmd/main.go
	./balance

run:
	go run ./cmd/main.go

down:
	migrate -path ./schema -database 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable' down

up:
	migrate -path ./schema -database 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable' up

build-docker:
	docker build -t balance-service .

compose:
	docker-compose up --build balance-service

# use for development. inits postgres
psql-init:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

init-dev:
	make psql-init
	make build

test:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
