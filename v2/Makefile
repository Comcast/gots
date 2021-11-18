DEFAULT: test

.PHONY: test
test: 
	go test -race `go list ./... | grep -v vendor`
