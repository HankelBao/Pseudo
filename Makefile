ifeq ($(OS),Windows_NT)
    uname_S := Windows
else
    uname_S := $(shell uname -s)
endif

default:
	@go run ./cmd/pseudo/ ./test/test.pse