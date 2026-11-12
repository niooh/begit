.PHONY: build test clean

build:
	go build -o bin/begit ./cmd/begit

test:
	go test ./...

clean:
	rm -rf bin
