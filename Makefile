build:
	@go build -o bin/finance cmd/main.go

run: build
	@./bin/finance
