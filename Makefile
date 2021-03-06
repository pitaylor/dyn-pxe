.PHONY: build clean start

build: out/dyn-pxe-darwin-amd64 out/dyn-pxe-linux-amd64

out/dyn-pxe-darwin-amd64: $(wildcard *.go)
	GOOS=darwin GOARCH=amd64 go build -o $@ $^

out/dyn-pxe-linux-amd64: $(wildcard *.go)
	GOOS=linux GOARCH=amd64 go build -o $@ $^

start:
	go run .

clean:
	rm -rf out
