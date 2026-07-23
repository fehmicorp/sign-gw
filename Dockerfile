FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o sign-gw \
    ./cmd/server
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache \
    ca-certificates \
    tzdata
COPY --from=builder /app/sign-gw .
COPY config.yaml .
RUN mkdir -p \
    /app/data \
    /app/data/eml \
    /app/log \
    /app/data/template

VOLUME ["/app/data"]
VOLUME ["/app/log"]

EXPOSE 25/tcp
ENV TZ=Asia/Kolkata
ENTRYPOINT ["./sign-gw"]