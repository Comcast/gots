DEFAULT: test

.PHONY: test
test:
	@echo "[-] Running UnitTests..."
	go test -v -short -race -parallel 2 ./...
