#! /bin/sh
set -ex

# =============================================================
# Removes volume from Docker image:
# An input base image is pulled from Docker Hub which provides
# a volume we need to get rid of in order to have data directly
# written to the container file system. The docker-copyedit.py
# metadata editor copies the base image, removes the VOLUME
# instruction from that copy's metadata and saves the copy as
# a new image with a given NAME:TAG. Finally, the base image
# is removed.
#
# Requires:
#  docker
#  wget
#  python
#
# USE ONLY in a CLEAN ENVIRONMENT.
# DO NOT USE this script if the base image is already found
# in your local registry or the new image tag belongs to an
# image in your local registry.
# =============================================================

# $REPO for $BASE image
REPO=mdillon/postgis

# $DIGEST for $BASE image
# Run 'docker pull [imgname]' first to get the $DIGEST of the
# target image. Then run 'docker image rm [imgname]'
DIGEST=sha256:c9eca7d32159529fd191239d74f2d11cb5dac1342eab426b79493bf47f933c9b

# Temporary tag for $BASE image
BASE="$REPO":tmp

# Docker volume to be removed in $BASE image
VOLUME=/var/lib/postgresql/data

# New image name after volume removal
NEW=siose-innova/postgis:2.5.2

# Pull image from Docker Hub
docker pull "$REPO"@"$DIGEST"

# Tag pulled image
docker tag "$REPO"@"$DIGEST" "$BASE"

# Get docker image metadata editor
wget https://raw.githubusercontent.com/gdraheim/docker-copyedit/master/docker-copyedit.py -O cmd.py
chmod +x cmd.py

# Remove $VOLUME from $BASE image and output into $NEW image
./cmd.py FROM "$BASE" INTO "$NEW" -vv REMOVE VOLUME "$VOLUME" 

# Remove temp directory
rm -rf load.tmp

# Remove $BASE image
docker image rm "$BASE"
