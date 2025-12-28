.PHONY: build test install clean

BINARY=jt

build:
	go build -o $(BINARY)

test:
	go test ./...

install:
	go install

clean:
	rm -f $(BINARY)
