BINARY := block_chain

all: build run

build:
	@echo "==> go build"
	@go build -o $(BINARY)

run:
	@echo "==> running"
	@./$(BINARY)

.PHONY: build run