FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
# Optional: use a mirror in networks where the default proxy is slow.
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
COPY assets/default-avatar.svg /app/defaults/avatar/default.svg
COPY scripts/start-http.sh /app/start-http.sh
COPY scripts/start-listener.sh /app/start-listener.sh

RUN chmod +x /app/start-http.sh /app/start-listener.sh \
  && sed -i 's/\r$//' /app/start-http.sh /app/start-listener.sh \
  && mkdir -p /app/storage/avatar /app/storage/cover /app/storage/video /app/defaults/avatar \
  && cp /app/defaults/avatar/default.svg /app/storage/avatar/default.svg

EXPOSE 8081

CMD ["/app/start-http.sh"]
