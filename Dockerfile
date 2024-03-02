FROM golang:1.22-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY .env.local .

RUN apk add --no-cache gcc musl-dev make
RUN go install github.com/cosmtrek/air@latest

