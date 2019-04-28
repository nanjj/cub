all: test

tidy:
	go mod tidy

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
