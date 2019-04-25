all: test

tidy:
	go mod tidy

test:
	go test -v ./...
