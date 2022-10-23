# syntax=docker/dockerfile:1

FROM golang:1.19.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy

COPY . ./

RUN go build -o /mail-service ./cmd/mail-service

ENTRYPOINT [ "/mail-service" ]