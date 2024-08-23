FROM golang:1.22 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct \
    CGO_ENABLED=0

WORKDIR /builder
COPY . .
RUN go mod download && \
    go build -o MediaWarp

FROM alpine:latest
COPY --from=builder /builder/MediaWarp /MediaWarp

RUN chmod +x /MediaWarp

EXPOSE 9000
VOLUME ["/config", "/logs"]
ENTRYPOINT ["/MediaWarp"]
