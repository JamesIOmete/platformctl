BIN=./bin/platformctl
PKG=./cmd/platformctl

.PHONY: build test ua clean

build:
	go build -o $(BIN) $(PKG)

test:
	go test ./...

ua: build
	./scripts/ua.sh

clean:
	rm -f $(BIN)
