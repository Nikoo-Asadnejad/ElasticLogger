FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o elasticlogger main.go

EXPOSE 8080

CMD ["./elasticlogger"]