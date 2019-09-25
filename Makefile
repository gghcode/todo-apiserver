.EXPORT_ALL_VARIABLES:
TEST_POSTGRES_DRIVER=postgres
TEST_POSTGRES_HOST=127.0.0.1
TEST_POSTGRES_PORT=5431
TEST_POSTGRES_USER=postgres
TEST_POSTGRES_NAME=postgres
TEST_POSTGRES_PASSWORD=postgres

dependency:
	@go get -v ./...
	@go get -u github.com/swaggo/swag/cmd/swag
	@go get -u github.com/oxequa/realize

live:
	@realize start --run --fmt --no-config

unit:
	@go test -race -v -short ./...

unit_ci:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -short ./...

integration: docker_up
	@go test -race -v -run Integration ./... || $(MAKE) docker_down
	@$(MAKE) docker_down

integration_ci: docker_up
	@TEST_POSTGRES_HOST=$$(ifconfig eth0 | grep "inet addr" | cut -d ':' -f 2 | cut -d ' ' -f 1) go test -race -coverprofile=coverage.txt -covermode=atomic -v -run Integration ./... || ($(MAKE) docker_down && exit 1)
	@$(MAKE) docker_down

docker_up: docker_down
	@-docker-compose -p integration -f docker-compose.integration.yml up -d

docker_down:
	@docker-compose -p integration -f docker-compose.integration.yml down -v

postgres:
	@docker run -d --name apas_postgres -p 5432:5432 postgres:11.3-alpine
 
redis:
	@docker run -d --name apas_redis -p 6379:6379 redis:5.0.5-alpine
