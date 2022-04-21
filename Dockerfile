FROM golang:alpine as builder

WORKDIR /spiderServer
COPY . .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache --virtual .build-deps wget git tar \
    && apk add --no-cache gdb binutils libc6-compat \
    && apk add --no-cache tzdata \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o spider .

FROM alpine:latest

WORKDIR /spiderServer

COPY --from=builder /spiderServer ./
COPY --from=builder /spiderServer/config.docker.yaml ./
EXPOSE 9000
ENTRYPOINT ./spider -c config.docker.yaml