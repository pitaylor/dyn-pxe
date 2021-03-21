#!/usr/bin/env bash

set -eo pipefail
set -x

host=${1:?host is missing}
tmp_dir=/tmp/dyn-pxe.tmp

ssh "${host}" -- bash -c "'mkdir -p ${tmp_dir} && rm -rf ${tmp_dir}/*'"
scp out/dyn-pxe-linux-amd64 config/dyn-pxe*.{service,path} "${host}:${tmp_dir}"
ssh -q "${host}" <<-SHELL
  set -x
  sudo systemctl stop dyn-pxe.service dyn-pxe-watcher.path || echo ok;
  sudo cp ${tmp_dir}/dyn-pxe-linux-amd64 /usr/local/bin/dyn-pxe \
    && sudo cp ${tmp_dir}/dyn-pxe*.{service,path} /etc/systemd/system \
    && sudo systemctl daemon-reload \
    && sudo systemctl enable dyn-pxe-watcher.{path,service} \
    && sudo systemctl start dyn-pxe.service dyn-pxe-watcher.path
SHELL

set +x
echo "Done!"
