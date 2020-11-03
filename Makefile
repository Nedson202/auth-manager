build_push:
	@echo "=============Building docker image============="
	docker build -f prod.Dockerfile -t samsonnegedu/auth-manager-api:1.0.0 .
	@echo "=============Pushing docker image to docker hub============="
	docker push samsonnegedu/auth-manager-api:1.0.0

logs:
	docker-compose logs -f

start:
	@echo "=============Starting api in docker============="
	docker-compose -f docker-compose.prod.yml up -d --build
	@echo "=============Loading service logs============="
	docker-compose -f docker-compose.prod.yml logs -f

start_dev:
	@echo "=============Starting api in docker============="
	docker-compose -f docker-compose.dev.yml up -d --build

	@echo "=============Loading service logs============="
	docker-compose -f docker-compose.dev.yml logs -f

stop:
	docker-compose -f docker-compose.dev.yml down

test:
	go test -v -cover ./...

dev_local:
	@echo "=============Starting api in development mode============="
	compileDaemon -build="go build -o bin/auth-manager ." -command="./bin/auth-manager" -color -graceful-kill
