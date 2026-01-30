.PHONY: build test install clean

build:
	go build -o bin/gitflow ./cmd/gitflow

test:
	go test ./...

install:
	go install ./cmd/gitflow

clean:
	rm -rf bin
