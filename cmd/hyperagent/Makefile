all: hyperagent

test: generate
	go test -v ./...

hyperagent: generate
	go build

generate:
	go generate
