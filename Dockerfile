# -------------------------------------------------------
# Build Stage
# -------------------------------------------------------

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

# -------------------------------------------------------
# Runtime Stage
# -------------------------------------------------------

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    tzdata\
    opensslx

COPY --from=builder /app/sign-gw .

RUN mkdir -p \
    /app/data \
    /app/data/eml \
    /app/data/templates \
    /app/data/certs \
    /app/logs

VOLUME ["/app/data"]
VOLUME ["/app/logs"]

EXPOSE 25
EXPOSE 587

ENV TZ=Asia/Kolkata

ENTRYPOINT ["./sign-gw","-c","/app/config.yaml"]