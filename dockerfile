FROM golang:1.20 as builder

COPY . /app
WORKDIR /app

RUN go get

EXPOSE 8080

ENV DB_HOST=host.docker.internal
ENV DB_PORT=3306
ENV DB_NAME=muchat
ENV DB_USER=ubuntu
ENV DB_PASS=ubuntu

# 等待 mysql 完全启动
ENV WAIT_VERSION 2.11.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD go build -o bin/go_another_chatgpt && ./bin/go_another_chatgpt
