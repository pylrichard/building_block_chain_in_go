BINARY := block_chain.exe

all: build test

build: deps
	@echo "==> go build"
	@go build -o $(BINARY)

deps:
	@go get -u go.etcd.io/bbolt

test:
	./$(BINARY) print_chain
	./$(BINARY) add_block -data "Send 1 BTC to Ivan"
	./$(BINARY) add_block -data "Pay 0.33 BTC for a coffee"
	./$(BINARY) print_chain

clean:
	rm block_chain.exe block_chain.db

.PHONY: build deps test