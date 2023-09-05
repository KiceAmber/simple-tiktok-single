# 设置基础镜像，用于编译代码
FROM golang:alpine AS builder

# 为镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o myapp .

# 新建一个新的小镜像用于运行服务
FROM debian:buster-slim

COPY ./wait-for.sh /
COPY ./manifest/ /manifest

COPY --from=builder /build/myapp /

RUN set -eux; \
    apt-get update; \
    apt-get install -y \
        --no-install-recommends \
        netcat; \
        chmod 755 wait-for.sh

EXPOSE 8989

# 设置容器启动命令
# ENTRYPOINT ./wait-for.sh mysql899:3306 redis300:6379 -- ./simple_tiktok
