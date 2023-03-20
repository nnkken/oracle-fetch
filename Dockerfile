FROM golang:1.20-alpine AS base

RUN apk update && apk add --no-cache build-base git bash curl linux-headers ca-certificates
RUN mkdir -p ./app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base

COPY . .
RUN go build -o /bin/oracle-fetch main.go
WORKDIR /bin
