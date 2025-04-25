hello:
	@echo "Hello, World!"

build:
	go build -o ./bin/TitanChain

run:build
	./bin/TitanChain

test:
	go test -v ./...