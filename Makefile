dependency:
	@go get -v ./...
	@go get -u github.com/swaggo/swag/cmd/swag
	@go get -u github.com/oxequa/realize

build:
	@go build

live:
	@realize start --run --fmt --no-config

run:
	@go run .

unit:
	@go test -race -v -short ./...

unit_ci:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -short ./...
