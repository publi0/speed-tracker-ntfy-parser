FROM alpine:latest

WORKDIR /root/

COPY ntfy-parser .

CMD ["./ntfy-parser"]