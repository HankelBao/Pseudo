ifeq ($(OS),Windows_NT)
    uname_S := Windows
else
    uname_S := $(shell uname -s)
endif

default:
ifeq ($(uname_S), Windows)
	@go run ./cmd/pseudo/ ./test/test.pse
	@clang ./tmp/test.ll -o ./tmp/test.exe
	@./tmp/test.exe
endif

ifeq ($(uname_S), Linux)
	@go run ./cmd/pseudo/ ./test/test.pse
	@clang ./tmp/test.ll -o ./tmp/test
	@./tmp/test
endif