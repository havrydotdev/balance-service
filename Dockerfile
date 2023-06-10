# BUILD STAGE
FROM golang:alpine AS builder

WORKDIR /build

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

# install migare cli so we can use it in prod stage
RUN GOBIN=/usr/local/bin/ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# build 
RUN go build -o balance cmd/main.go

# PROD STAGE
# use alpine to reduce image`s size
FROM alpine

WORKDIR /build

# copy built exe file
COPY --from=builder /build/balance /build/balance

# copy wait-for-postgres.sh
COPY --from=builder /build/scripts /build/scripts

# copy .env file
COPY --from=builder /build/.env /build/.env

# copy migrations dir from build
COPY --from=builder /build/schema /build/schema

# copy golang-migrate
COPY --from=builder /usr/local/bin/ /usr/local/bin/

# copy configs
COPY --from=builder /build/configs /build/configs

# install postgresql-client
RUN apk update
RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x ./scripts/wait-for-postgres.sh

# run service
CMD ["./balance"]