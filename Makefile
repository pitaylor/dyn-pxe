.PHONY: build clean deploy pre_deploy start

build: out/lab-init-darwin-amd64 out/lab-init-linux-amd64

out/lab-init-darwin-amd64: $(wildcard *.go)
	GOOS=darwin GOARCH=amd64 go build -o $@ $^

out/lab-init-linux-amd64: $(wildcard *.go)
	GOOS=linux GOARCH=amd64 go build -o $@ $^

start:
	go run *.go -dataDir . -templateDir ./templates -execDir ./exec

define DEPLOY_SCRIPT
set -e
# Install binary
mv -f $(DEPLOY_TMP)/lab-init /usr/local/bin
chown root.root /usr/local/bin/lab-init
# Install templates
mkdir -p /etc/lab-init/templates
cp -rn $(DEPLOY_TMP)/templates/* /etc/lab-init/templates
# Install systemd script
cp -f $(DEPLOY_TMP)/lab-init.service /etc/systemd/system
systemctl start lab-init
systemctl enable lab-init
# Cleanup
rm -rf $(DEPLOY_TMP)
endef

DEPLOY_TMP=/tmp/lab-init.tmp

export DEPLOY_SCRIPT

deploy: build
ifndef HOST
	$(error HOST is not set)
endif
	ssh $(HOST) -- rm -rf $(DEPLOY_TMP)
	ssh $(HOST) -- mkdir -p $(DEPLOY_TMP)/templates
	scp -q out/lab-init-linux-amd64 $(HOST):$(DEPLOY_TMP)/lab-init
	scp -q templates/* $(HOST):$(DEPLOY_TMP)/templates
	scp -q config/lab-init.service $(HOST):$(DEPLOY_TMP)
	echo "$$DEPLOY_SCRIPT" | ssh $(HOST) -- "sudo bash"
	@echo "Done."

clean:
	rm -rf out
