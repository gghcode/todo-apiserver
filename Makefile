dependency:
	@go get -v ./...

build:
	@go build

run:
	@go run .

unit: dependency
	@go test -race -v -short ./...

unit_ci: dependency
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -short ./...
