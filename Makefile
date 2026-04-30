.PHONY: run test tidy build

run:
	go run main.go

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o goreview
