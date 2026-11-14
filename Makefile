.PHONY: test clean path install

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

release:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/begit ./cmd/begit
