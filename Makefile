.PHONY: build clean deploy pre_deploy start

build: out/lab-init-darwin-amd64 out/lab-init-linux-amd64

out/lab-init-darwin-amd64: *.go
	GOOS=darwin GOARCH=amd64 go build -o $@ $^

out/lab-init-linux-amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o $@ $^

start:
	go run *.go -dataDir . -templateDir ./templates -execCmd ./scripts/exec.sh

define DEPLOY_SCRIPT
set -euo pipefail
mv -f /tmp/lab-init.tmp /usr/local/bin/lab-init
chmod 755 /usr/local/bin/lab-init
mkdir -p /etc/lab-init
cp -n /tmp/lab-init/* /etc/lab-init
rm -rf /tmp/lab-init
endef

export DEPLOY_SCRIPT

pre_deploy:
ifndef HOST
	$(error HOST is not set)
endif

deploy: pre_deploy build
	ssh $(HOST) -- mkdir -p /tmp/lab-init
	scp -q out/lab-init-linux-amd64 $(HOST):/tmp/lab-init.tmp
	scp -q templates/* scripts/exec.sh $(HOST):/tmp/lab-init
	echo "$$DEPLOY_SCRIPT" | ssh $(HOST) -- "sudo bash"
	echo "Done."

clean:
	rm -rf out
