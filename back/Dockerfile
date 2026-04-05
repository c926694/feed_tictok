FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
# 在 go mod download 前设置代理
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/http ./cmd/http
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/listener ./cmd/listener

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/http /app/http
COPY --from=builder /out/listener /app/listener
COPY config/config.compose.example.yaml /app/config/config.yaml

RUN mkdir -p /app/storage/avatar /app/storage/cover /app/storage/video

EXPOSE 8080

CMD ["/app/http"]

