FROM golang:1.22-alpine AS builder

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
COPY config.yaml .

RUN mkdir -p /app/data

EXPOSE 8080

ENTRYPOINT ["./opstimer"]
CMD ["--config", "config.yaml"]
