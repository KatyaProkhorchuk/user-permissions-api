
FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go mod tidy
EXPOSE ${PORT_CLIENT}
EXPOSE ${PORT_ADMIN}

RUN mkdir logs && touch ./logs/service_logs.log

VOLUME ./logs
CMD go run ./cmd/main.go
