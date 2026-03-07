# Build Stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache upx

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the statically linked binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o go-samba4 .
RUN upx --best --lzma go-samba4

# Runtime Stage
FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata sqlite-libs

COPY --from=builder /app/go-samba4 .
COPY config.toml /etc/go-samba4/config.toml

EXPOSE 8080

CMD ["./go-samba4", "serve", "--config", "/etc/go-samba4/config.toml"]
