# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . ./

RUN go build -o /weather-bot main.go

EXPOSE 8080

CMD [ "/weather-bot" ]