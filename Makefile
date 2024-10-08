build:
	@go build -o bin/finance cmd/main.go

test:
	@test_packages=$$(./scripts/filter_test_packages.sh); \
	go test -v $$test_packages

run: build
	@./bin/finance

gen:
	@go run github.com/99designs/gqlgen generate
	
