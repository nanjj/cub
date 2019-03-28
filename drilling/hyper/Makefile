all: hyper

test: generate
	go test -v ./...

hyper: generate
	go build

generate:
	go generate
