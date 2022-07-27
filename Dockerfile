FROM golang:alpine as builder

WORKDIR /stockInfo
COPY . .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache --virtual .build-deps wget git tar \
    && apk add --no-cache gdb binutils libc6-compat \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o stockApiServer .

FROM alpine:latest

WORKDIR /spiderServer

COPY --from=builder /stockInfo ./
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /stockInfo/config.docker.yaml ./
ENV TZ=Asia/Shanghai
EXPOSE 8000
ENTRYPOINT ./stockApiServer -c config.docker.yaml