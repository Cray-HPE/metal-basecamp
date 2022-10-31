#!/bin/bash
# MIT License
#
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
set -euf -o pipefail

if [ $# -lt 2 ]; then
    echo >&2 "usage: basecamp-init PIDFILE CIDFILE [CONTAINER [VOLUME]]"
    exit 1
fi

BASECAMP_PIDFILE="$1"
BASECAMP_CIDFILE="$2"
BASECAMP_CONTAINER_NAME="${3:-basecamp}"
BASECAMP_IMAGE_NAME="metal-basecamp"

BASECAMP_IMAGE_PATH="@@basecamp-path@@"
BASECAMP_IMAGE="@@basecamp-image@@"

BASECAMP_VOLUME_MOUNT_CONFIG='/var/www/ephemeral/configs:/app/configs:rw,shared'
BASECAMP_VOLUME_MOUNT_STATIC='/var/www/ephemeral/static:/app/static:rw,shared'

command -v podman >/dev/null 2>&1 || { echo >&2 "${0##*/}: command not found: podman"; exit 1; }

mkdir -pv "$(echo ${BASECAMP_VOLUME_MOUNT_STATIC} | cut -f 1 -d :)"
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


# Load basecamp image if it doesn't already exist
if ! podman image inspect "$BASECAMP_IMAGE" &>/dev/null; then
    # get the image id (there may be more than one in certain situations)
    BASECAMP_IMAGE_ID=()
    read -r -a BASECAMP_IMAGE_ID <<< "$(podman images --noheading --format "{{.Id}}" --filter label="org.label-schema.name=$BASECAMP_IMAGE_NAME")"

    # if more than one image is found, clear them out
    if [ ${#BASECAMP_IMAGE_ID[@]} -gt 1 ]; then
        podman image rm --force --ignore "${BASECAMP_IMAGE_ID[@]}"
    fi

    # otherwise, load and tag the image
    podman load -i "$BASECAMP_IMAGE_PATH" || exit
    # get the image ID
    read -r -a BASECAMP_IMAGE_ID <<< "$(podman images --noheading --format "{{.Id}}" --filter label="org.label-schema.name=$BASECAMP_IMAGE_NAME")"
    # tag the image appropriately
    podman tag "${BASECAMP_IMAGE_ID[@]}" "$BASECAMP_IMAGE"
fi

# always ensure fresh files
# since this runs in a container and called from a systemd service,
# the container could have been created with malformed json or some other error
# clearing it out everytime and making a fresh container allows for a clean start
podman stop --ignore -t 10 "$BASECAMP_CONTAINER_NAME"
podman rm --force --ignore "$BASECAMP_CONTAINER_NAME"

rm -f "$BASECAMP_PIDFILE" 
rm -f "$BASECAMP_CIDFILE"

# create a fresh container in case the previous failed to start due to malformed json
podman create \
    --conmon-pidfile "$BASECAMP_PIDFILE" \
    --cidfile "$BASECAMP_CIDFILE" \
    --cgroups=no-conmon \
    --net host \
    --volume $BASECAMP_VOLUME_MOUNT_CONFIG \
    --volume $BASECAMP_VOLUME_MOUNT_STATIC \
    --name "$BASECAMP_CONTAINER_NAME" \
    --env GIN_MODE="${GIN_MODE:-release}" \
    "$BASECAMP_IMAGE"

# start the container so we can get the exit code
podman start "$BASECAMP_CONTAINER_NAME"

# exit this script with the exit code from the container so it is passed appropriately to systemd
rc=""
rc=$(podman inspect --type container "${BASECAMP_CONTAINER_NAME}" --format "{{.State.ExitCode}}")
exit "${rc}"
