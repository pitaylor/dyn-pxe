.PHONY: build clean install-remote start test

build: out/dyn-pxe-darwin-amd64 out/dyn-pxe-linux-amd64

out/dyn-pxe-darwin-amd64: $(wildcard *.go)
	GOOS=darwin GOARCH=amd64 go build -o $@ $^

out/dyn-pxe-linux-amd64: $(wildcard *.go)
	GOOS=linux GOARCH=amd64 go build -o $@ $^

install-remote: build
	@scripts/install-remote.sh $(host)

start:
	go run . -config config.yml -static-dir static

test:
	go test -v ./...

clean:
	rm -rf out
