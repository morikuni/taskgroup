.PHONY: test
test:
	go test -count 1 -v -race ./...

mod:
	go mod tidy