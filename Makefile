build:
	@echo "=============building Local API============="
	docker build -f Dockerfile -t main .

logs:
	docker-compose logs -f

start:
	@echo "=============starting api locally============="
	docker-compose up -d
	@echo "=============Loading service logs============="
	make logs

stop:
	docker-compose down -v --rmi all

test:
	go test -v -cover ./...

dev:
	@echo "=============starting api in development mode============="
	compileDaemon -build="go build -o bin/user-service ." -command="./bin/user-service" -color -graceful-kill
