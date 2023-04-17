FROM golang:1.20 as builder

COPY . /app
WORKDIR /app

RUN go get

EXPOSE 8080

ENV DB_HOST=host.docker.internal
ENV DB_PORT=3306
ENV DB_NAME=chatgpt
ENV DB_USER=ubuntu
ENV DB_PASS=ubuntu

CMD go build -o bin/go_another_chatgpt && ./bin/go_another_chatgpt
