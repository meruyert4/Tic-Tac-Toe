FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o tic-tac-toe ./backend/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tic-tac-toe .
COPY frontend ./frontend

EXPOSE 8080

CMD ["./tic-tac-toe"]
