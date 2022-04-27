BINARY := block_chain.exe

all: build run

build:
	@echo "==> go build"
	@go build -o $(BINARY)

run:
	@echo "==> run"
	@./$(BINARY)

clean:
	rm block_chain.exe

.PHONY: build run