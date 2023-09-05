# 设置基础镜像
FROM golang:alpine AS build

# 设置拉取依赖的镜像源
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 将项目下的所有文件放到 /src 目录下
COPY . /src

# 进入 src 目录后编译项目生成可执行文件
RUN cd /src && go build -o myapp

FROM debian:stretch-slim

# 设置工作目录
WORKDIR /app

# 将二进制文件复制到工作目录 app 下
COPY --from=build /src/myapp /app/
COPY --from=build /src/manifest/ /app/manifest
COPY --from=build /src/wait-for.sh /app/

RUN apt-get update; \
    apt-get install -y \
        --no-install-recommends \
        netcat; \
        chmod 755 /app/wait-for.sh

EXPOSE 8989

# 设置容器启动命令
# ENTRYPOINT ./wait-for.sh mysql899:3306 redis300:6379 -- ./simple_tiktok
