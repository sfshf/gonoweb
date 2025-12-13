#!/usr/bin/env bash

# 开启本地环境 -- 数据库等服务
docker compose -f ./deploy/dev/compose.yml up -d

echo -e "\033[1;32mGenerate Gorm Models:\033[0m"
# 等待几秒，等前置环境都启动好
sleep 15
./script/gen/gorm.sh