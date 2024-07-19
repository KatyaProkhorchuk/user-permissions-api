FROM golang:1.22

WORKDIR /server

COPY . .

RUN go mod tidy

EXPOSE ${SERVER_PORT}

CMD ["go", "run", "./main.go"]