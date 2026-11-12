.PHONY: build test clean


build:
	go build -o bin/begit ./cmd/begit

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin

path:
	printf '\033]52;c;ZXhwb3J0IFBBVEg9IiRQQVRIOiQoZ28gZW52IEdPUEFUSCkvYmluIg==\a'

install:
	go install ./cmd/begit
