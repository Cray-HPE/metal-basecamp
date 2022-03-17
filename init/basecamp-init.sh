#!/bin/bash

set -x
if [ $# -lt 2 ]; then
    echo >&2 "usage: basecamp-init PIDFILE CIDFILE [CONTAINER [VOLUME]]"
    exit 1
fi

BASECAMP_PIDFILE="$1"
BASECAMP_CIDFILE="$2"
BASECAMP_CONTAINER_NAME="${3-basecamp}"
BASECAMP_VOLUME_NAME="${4:-${BASECAMP_CONTAINER_NAME}-configs}"

BASECAMP_IMAGE_PATH="@@basecamp-path@@"
BASECAMP_IMAGE="@@basecamp-image@@"

BASECAMP_VOLUME_MOUNT_CONFIG='/var/www/ephemeral/configs:/app/configs:rw,shared'
BASECAMP_VOLUME_MOUNT_STATIC='/var/www/ephemeral/static:/app/static:rw,shared'

command -v podman >/dev/null 2>&1 || { echo >&2 "${0##*/}: command not found: podman"; exit 1; }


# always ensure pid file is fresh
rm -f "$BASECAMP_PIDFILE"
mkdir -pv "$(echo ${BASECAMP_VOLUME_MOUNT_CONFIG} | cut -f 1 -d :)"
test -e "$(echo ${BASECAMP_VOLUME_MOUNT_CONFIG} | cut -f 1 -d :)/data.json" ||\
cat << EOF > "$(echo ${BASECAMP_VOLUME_MOUNT_CONFIG} | cut -f 1 -d :)/data.json"
{
  [
    // "mac": {metadata...}
  ]
}
EOF
# Set up a mutable, default file. Users reading this, can edit this or edit the
# actual created file. Editing here is persistent on restart.
test -e "$(echo ${BASECAMP_VOLUME_MOUNT_CONFIG} | cut -f 1 -d :)/server.yaml" ||\
cat << EOF > "$(echo ${BASECAMP_VOLUME_MOUNT_CONFIG} | cut -f 1 -d :)/server.yaml"
# Basecamp Configuration
bind: ":8888"
local-mode: true
local-data: "./configs/data.json"
serve-static: true
static-dir: "./static/"
EOF

mkdir -pv "$(echo ${BASECAMP_VOLUME_MOUNT_STATIC} | cut -f 1 -d :)"
# Create basecamp container
if ! podman inspect "$BASECAMP_CONTAINER_NAME" ; then
    rm -f "$BASECAMP_CIDFILE" || exit
    # Load basecamp image if it doesn't already exist
    if ! podman image inspect "$BASECAMP_IMAGE" >dev/null; then
        podman load "$BASECAMP_IMAGE_PATH" "$BASECAMP_IMAGE" || exit
    fi
    podman create \
        --conmon-pidfile "$BASECAMP_PIDFILE" \
        --cidfile "$BASECAMP_CIDFILE" \
        --cgroups=no-conmon \
        --net host \
        --volume $BASECAMP_VOLUME_MOUNT_CONFIG \
        --volume $BASECAMP_VOLUME_MOUNT_STATIC \
        --name "$BASECAMP_CONTAINER_NAME" \
        --env GIN_MODE="${GIN_MODE:-release}" \
        "$BASECAMP_IMAGE" || exit
    podman inspect "$BASECAMP_CONTAINER_NAME" || exit
fi
