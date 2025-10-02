FROM alpine:latest

COPY MediaWarp ./MediaWarp

RUN chmod +x ./MediaWarp

EXPOSE 9000
VOLUME ["/etc/localtime", "/etc/timezone", "/config", "/logs", "/custom"]
ENTRYPOINT ["/MediaWarp"]