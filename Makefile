.PHONY: build test vet fmt install clean

build:
	go build -ldflags "-s -w" -o shipyard .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

install: build
	install -m 0755 shipyard $(HOME)/.local/bin/shipyard

clean:
	rm -f shipyard
	rm -rf dist
