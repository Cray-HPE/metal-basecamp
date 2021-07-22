# Basecamp 

[![Continuous Integration](https://github.com/Cray-HPE/metal-basecamp/actions/workflows/ci.yml/badge.svg)](https://github.com/Cray-HPE/metal-basecamp/actions/workflows/ci.yml)


This repo provides the metadata datasource for cloud-init during bootstrap and recovery.

The sister-service, CRAY-BSS, provides this same capability to scale through Kubernetes. Baecamp can be
used in Docker or Podman quickly for bootstrapping a kubernetes cluster. The data served by Basecamp
is compatible with Cray-BSS, and is handed-off as part of the deployment process.

## Usage - Server

Basecamp can be ran multiple ways. For the full experience, we suggest using the Docker image.
Developers should use whichever environment they like, whether that's the Docker container
or their local Go-lang env.

### Podman (or Docker)

1. Create configuration directories

    ```bash
    $> mkdir -p configs
    $> touch configs/server.yaml # fill this in
    $> touch configs/data.json # fill this in
    ```
2. Run the container:

    ```bash
   $> image='arti.dev.cray.com/csm-docker-master-local/metal-basecamp'
   $> podman create --net host --volume $(pwd)/configs:/app/configs --name basecamp "$image"
    ```
3. You should now be able to run queries against the container:

    ```bash
    $> curl http://localhost:8888/meta-data
    $> curl http://localhost:8888/user-data
    ```

### Daemon

There is a systemd-daemon in this repo for running the Basecamp container in the background. The 
daemon can be obtained through linux package managers.

> The daemon will setup `configs/data.json` and `configs/server.yaml` if they do not already exist.

##### OpenSuSE / SLES

```bash
# Add repo. and install metal-basecamp
$> repo=http://car.dev.cray.com/artifactory/csm/MTL/sle15_sp2_ncn/
$> zypper addrepo --no-gpgcheck --refresh "$repo" csm-metal
$> zypper install metal-basecamp

# Enable and start
$> systemctl enable basecamp
$> systemctl start basecamp
```
### Logs

Basecamp outputs logs in two places; container logs and daemon logs.

```bash
# Container logs
$> podman logs -f basecamp

# Daemon logs
$> journalctl -xeu basecamp -f
``` 


## Usage - Client

Two steps for usage.

### 1. Install cloud-init

Your instance (virtual or metal) needs to have `cloud-init` installed.


##### OpenSuSE / SLES

```bash
$> zypper install cloud-init
$> systemctl enable cloud-configs cloud-init-local cloud-init cloud-final
```

### 2. Configure cloud-init

#### Option 1:

    # Add the following to `/etc/cloud/cloud.cfg`
    datasource:
      NoCloud:
        seedfrom: http://{IP OF METADATA SERVER}:{PORT}/

#### Option 2:

> Note: _This is what the Shasta Pre-Install Toolkit passes in its ipxe script._

##### iPXE
For `.ipxe` scripts, which escape the `;`.
```bash
# Add to /etc/default/grub CMDLINE_DEFAULT or to the kernel line for ipxe:
ds=nocloud-net;s=http://{IP OF METADATA SERVER}:{PORT}/
```

##### GRUB
For `grub.cfg` where the `;` is not escaped by default (so we need a back-slash).
```bash
# For `grub.cfg`
ds=nocloud-net\;s=http://{IP OF METADATA SERVER}:{PORT}/
```
