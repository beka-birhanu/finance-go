build:
	@go build -o bin/finance cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/finance
