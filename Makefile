.EXPORT_ALL_VARIABLES:
PORT=8080
BASE_PATH=api
BCRYPT_COST=1
GRACEFUL_SHUTDOWN_TIMEOUT_SEC=30
JWT_SECRET_KEY=debug-key
JWT_ACCESS_EXPIRES_IN_SEC=3600
JWT_REFRESH_EXPIRES_IN_SEC=86400
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5431
POSTGRES_USER=postgres
POSTGRES_NAME=postgres
POSTGRES_PASSWORD=postgres
REDIS_ADDR=127.0.0.1:6378
REDIS_PASSWORD=test

run:
	@go run .

live:
	@realize start --run --fmt --no-config

unit:
	@go test -race -short ./...

unit_ci:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -short ./...

integration: docker_up
	@go test -race -run Integration ./... || $(MAKE) docker_down
	@$(MAKE) docker_down

integration_ci: docker_up
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -run Integration ./... || ($(MAKE) docker_down && exit 1)
	@$(MAKE) docker_down

docker_up: docker_down
	@-docker-compose -p integration -f docker-compose.integration.yml up -d

docker_down:
	@docker-compose -p integration -f docker-compose.integration.yml down -v
