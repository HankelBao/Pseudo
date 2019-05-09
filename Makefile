linux:
	@go run ./cmd/gopse/ ./test/test.pse
	@clang ./tmp/test.ll -o ./tmp/test
	@./tmp/test

windows:
	@go run ./cmd/gopse/ ./test/test.pse
	@clang ./tmp/test.ll -o ./tmp/test.exe
	@./tmp/test.exe
