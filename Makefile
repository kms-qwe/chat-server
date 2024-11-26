include ./build/.env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=${MIGRATION_DIR}
LOCAL_MIDGRATION_DSN="host=localhost port=${PG_PORT} dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"


install-golangci-lint:
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	cd app && ${LOCAL_BIN}/golangci-lint run ./... --config ../.golangci.pipeline.yaml

install-deps:
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=${LOCAL_BIN} go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=${LOCAL_BIN} go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=${LOCAL_BIN} go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.2

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


local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_MIGRATION_DSN} down -v


local-up:
	docker compose -f ./build/docker-compose.local.yaml up --build -d 

local-down:
	docker compose -f ./build/docker-compose.local.yaml down 

test:
	cd app && go clean -testcache
	cd app && go test ./... -covermode count -coverpkg=./internal/service/...,./internal/api/... -count 5

test-coverage:
	cd app && go clean -testcache
	cd app && go test ./... -coverprofile=../coverage.tmp.out -covermode count -coverpkg=./internal/service/...,./internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	cd app && go tool cover -html=../coverage.out;
	cd app && go tool cover -func=../coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore