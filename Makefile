BINARY_PATH=build
BINARY_NAME=app
SRC=cmd/main.go

.PHONY: all build run test clean help

all: build

build:
	go build -o $(BINARY_PATH)/$(BINARY_NAME) $(SRC)

run: build
	$(BINARY_PATH)/$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -rf $(BINARY_PATH)