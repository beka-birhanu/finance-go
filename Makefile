build:
	@go build -o bin/finance cmp/main.go

run: build
	@./bin/finance
