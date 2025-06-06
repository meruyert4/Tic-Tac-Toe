FROM golang:1.23

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o tic-tac-toe .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tic-tac-toe .
COPY frontend/ ./frontend/

EXPOSE 8080

CMD ["./tic-tac-toe"]