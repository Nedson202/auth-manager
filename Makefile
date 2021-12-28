#!make
all:
	$(error Pick a target.)

.PHONY: proto_gen

logs:
	docker-compose logs -f

test:
	go test -v -cover ./...

proto_gen:
	@echo "=============Generating proto============="
	mkdir -p api/proto/v1 && protoc -I ./api/proto/v1 \
		-I. -I$$GOPATH/src/github.com/nedson202/auth-manager/third_party \
		--go_out ./api/proto/v1 --go_opt=paths=source_relative \
		--validate_out=paths=source_relative,"lang=go:./api/proto/v1" \
		--grpc-gateway_out=logtostderr=true:api/proto/v1 --grpc-gateway_opt paths=source_relative \
		api/proto/v1/auth.proto

swagger_gen:
	@echo "=============Generating openapi spec============="
	mkdir -p api/swagger/v1 && protoc -I ./api/proto/v1 \
		-I. -I$$GOPATH/src/github.com/nedson202/auth-manager/third_party \
		--go-grpc_out ./api/proto/v1 --go-grpc_opt paths=source_relative \
		--swagger_out=logtostderr=true:api/swagger/v1 \
		api/proto/v1/auth.proto

start: proto_gen
	@echo "=============Starting api in docker============="
	docker-compose -f docker-compose.prod.yml up -d --build
	@echo "=============Loading service logs============="
	docker-compose -f docker-compose.prod.yml logs -f

dev_local: proto_gen
	@echo "=============Starting api in development mode============="
	cd cmd/server && compileDaemon -build="go build -o bin/auth-manager-rpc ." -command="./bin/auth-manager-rpc" -color -graceful-kill

start_dev:
	@echo "=============Starting api in docker============="
	docker-compose -f docker-compose.dev.yml up -d --build

	@echo "=============Loading service logs============="
	docker-compose -f docker-compose.dev.yml logs -f

build_push: proto_gen
	@echo "=============Building docker image============="
	docker build -f prod.Dockerfile -t samsonnegedu/auth-api:1.0.2 .
	@echo "=============Pushing docker image to docker hub============="
	docker push samsonnegedu/auth-api:1.0.2

load_test:
	artillery run load_testing/load-test.yaml
