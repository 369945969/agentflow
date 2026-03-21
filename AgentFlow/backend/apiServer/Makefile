.PHONY: build run test clean

BINARY_NAME=server

build:
	go build -o $(BINARY_NAME) main.go

run: build
	./$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f server.log
