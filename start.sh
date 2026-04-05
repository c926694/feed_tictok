#!/bin/bash
# ===============================
# start.sh - 启动 feed_tictok 项目 Docker 容器
# 假设项目源码已经存在
# ===============================

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose CLI 插件
if ! docker compose version &> /dev/null; then
    echo "Docker Compose CLI 插件未安装或不可用"
    exit 1
fi

DOCKER_COMPOSE_CMD="docker compose"

# 确保在项目根目录（包含 docker-compose.yml）
if [ ! -f docker-compose.yml ]; then
    echo "docker-compose.yml 文件不存在，请确认在项目根目录运行脚本"
    exit 1
fi

# 构建并启动 Docker 容器
echo "开始构建 Docker 镜像并启动容器..."
$DOCKER_COMPOSE_CMD up -d --build || { echo "Docker Compose 启动失败"; exit 1; }

# 打印容器状态
echo "---------------------------------"
echo "容器状态："
$DOCKER_COMPOSE_CMD ps
echo "---------------------------------"

# 获取 Linux IP 地址
LINUX_IP=$(hostname -I | awk '{print $1}')

echo "前端 Vue 页面访问地址： http://$LINUX_IP/"
echo "后端 API 访问地址：     http://$LINUX_IP:8080/api/"
echo "RabbitMQ 管理界面：     http://$LINUX_IP:15672 (默认 guest/guest)"
echo "---------------------------------"
echo "容器启动完成 ✅"