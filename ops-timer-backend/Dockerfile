FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/opstimer ./cmd/server/

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/opstimer .

RUN mkdir -p /app/data

EXPOSE 8080

# 默认查找 ./config.yaml；镜像内可不打包配置文件，通过 TIMER_* 环境变量注入（见 README）
ENTRYPOINT ["./opstimer"]
