# =========================
# BUILD STAGE
# =========================
FROM golang:1.24-alpine AS builder

# install git for GOPROXY=direct
RUN apk add --no-cache git

WORKDIR /app

# better caching
COPY go.mod go.sum ./
RUN GOPROXY=direct go mod download

COPY . .

RUN go build -o main ./cmd/server

# =========================
# RUNTIME STAGE
# =========================
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]