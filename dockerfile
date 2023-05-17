FROM golang:1.20 as builder

COPY . /app
WORKDIR /app

RUN go get

RUN go build -o bin/go_another_chatgpt
RUN sh dk_init.sh

# 等待 mysql 完全启动
ENV WAIT_VERSION 2.11.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

FROM golang:1.20

EXPOSE 8080

ENV DB_HOST=host.docker.internal
ENV DB_PORT=3306
ENV DB_NAME=muchat
ENV DB_USER=ubuntu
ENV DB_PASS=ubuntu

WORKDIR /app

COPY --from=builder /wait /wait
COPY --from=builder /app/bin/go_another_chatgpt /app/bin/go_another_chatgpt
COPY --from=builder /app/配置 /app/配置
COPY --from=builder /app/storage /app/storage
COPY --from=builder /app/resources /app/resources


CMD /app/bin/go_another_chatgpt