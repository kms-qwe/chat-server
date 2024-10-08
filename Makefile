LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p app/pkg/chat_v1
	protoc --proto_path app/api/chat_v1 \
	--go_out=app/pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=app/pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	app/api/chat_v1/chat.proto


build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/service_linux ./app/cmd/grpc_server/main.go

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o ./bin/service_static_linux ./app/cmd/grpc_server/main.go

copy-to-server:
	scp ./bin/service_linux root@185.104.251.226:~/Krasnobaev/grpc_chat-server

copy-static-to-server:
	scp ./bin/service_static_linux root@185.104.251.226:~/Krasnobaev/grpc_chat-server

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/olezhek28/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAA6ELvdoGnA7EL4qSbkICKeDkVldDC2OOU cr.selcloud.ru/olezhek28
	docker push cr.selcloud.ru/olezhek28/test-server:v0.0.1
	