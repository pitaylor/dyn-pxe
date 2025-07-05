# Dynamic PXE Boot Service

A PXE boot server for provisioning homelab machines. Serves dynamic content over HTTP and TFTP.

## What it does

- Serves content via HTTP (port 8080) and TFTP (port 69) for PXE booting
- Generates dynamic content using templates or by running commands
- Captures URL path parameters and passes them to templates/commands
- Falls back to serving static files
- Configured with YAML files

## Example Uses

### PXE Boot Menu
Create dynamic boot menus that change based on the machine's MAC address:

```yaml
# config.yml

# Define HTTP routes that map to go templates. Routes can include named 
# parameters that are passed to templates.
resources:
  - route: pxelinux.cfg/default
    template: pxelinux.cfg
  - route: pxelinux.cfg/01-(?P<MacAddress>[\da-f]{2}(-[\da-f]{2}){5})
    template: pxelinux.cfg

# Define variables that can be used by go templates
variables:
  BaseUrl: http://192.168.1.100:8069
  Hosts:
    f4-4d-30-61-9a-6c: lab1
    f4-4d-30-61-cd-39: lab2
```

The template can use the captured MAC address to customize the boot menu:
```
{{ $hostname := or (index .Vars.Hosts .Params.MacAddress) .Params.MacAddress }}
DEFAULT menu.c32
MENU TITLE Boot Menu - {{ $hostname }}

LABEL ubuntu-22.04-pre
  MENU LABEL Install Ubuntu 22.04 (Preconfigured)
  KERNEL ubuntu-22.04.1/vmlinuz
  INITRD ubuntu-22.04.1/initrd
  APPEND autoinstall ds=nocloud-net;s={{ .Vars.BaseUrl }}/cloud-init/{{ .Params.MacAddress }}/
```

### Cloud-Init Provisioning
Serve customized cloud-init configs for automated Ubuntu installs:

```yaml
resources:
  - route: cloud-init/(?P<MacAddress>[\da-f]{2}(-[\da-f]{2}){5})/user-data
    template: user-data.yml
```

The template generates hostname and configuration based on MAC address:
```yaml
{{ $hostname := or (index .Vars.Hosts .Params.MacAddress) (.Params.MacAddress | shortHash | printf "lab-%s") }}
autoinstall:
  identity:
    hostname: {{ $hostname }}
    username: example-username
```

## Remote Install

```bash
make install-remote host=<target-host>
```
