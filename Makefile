.PHONY: run test lint

run:
	go run cmd/tt4em/main.go

test:
	go test -race ./...

lint:
	golangci-lint run
