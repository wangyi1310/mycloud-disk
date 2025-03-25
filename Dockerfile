# mycloud-disk 编译环境
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mycloud-disk .

# mycloud-disk 运行环境
FROM alpine:latest
WORKDIR /usr/local/services
COPY --from=builder /app/mycloud-disk .

# 安装时区配置，设置权限
RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x ./ mycloud-disk 
CMD ["./mycloud-disk"]