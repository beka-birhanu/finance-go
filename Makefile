build:
	@go build -o bin/finance api/main.go

run: build
	@./bin/finance
