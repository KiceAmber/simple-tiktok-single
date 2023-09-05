#!/bin/bash

# 检测 docker 环境
if ! command -v docker &> /dev/null; then
  echo "Error: docker command not found"
  exit 1
fi

# 检测 docker-compose j环境
if ! command -v docker-compose &> /dev/null; then
  echo "Error: docker-compose command not found"
  exit 1
fi

# 拉取 mysql 和 redis 镜像
docker pull mysql:8.0.19
docker pull redis:5.0.7

# 启动服务
docker-compose up -d
