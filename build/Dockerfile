FROM golang:1.22-alpine3.19 AS builder

COPY app /github.com/kms-qwe/chat-server/source/
WORKDIR /github.com/kms-qwe/chat-server/source/

RUN go mod download
RUN go build -o ./bin/crud_server ./cmd/main.go

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /github.com/kms-qwe/chat-server/source/bin/crud_server .

ARG PATH_TO_CONFIG_FILE
COPY ${PATH_TO_CONFIG_FILE} "./config/.env"

CMD ["./crud_server", "-config-path=./config/.env"]