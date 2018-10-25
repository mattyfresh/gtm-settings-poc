#!/bin/bash

# @TODO - this is a WIP!!!  Still need to get the correct AZ values to publish this

# This script deploys the latest version of the gobot container image to
# a free-floating container hosted by Azure Container Instances

# if this job was called by the previous job, as part of the normal workflow
# triggered by the push to the repo, we expect to find the file with the docker image
# name to be deployed in docker_image_full_name.txt

set -e



if [ -f docker_image_full_name.txt ]; then
  echo "Found docker_image_full_name.txt"
  DOCKER_IMAGE_FULL_NAME=$(cat docker_image_full_name.txt)
fi

if [ -f version.txt ]; then
  echo "Found version.txt"
  APP_VERSION=$(cat version.txt)
fi

# if the job was triggered using the API, the docker image name
# must be passed as a parameter
if [ "$DOCKER_IMAGE_FULL_NAME" = "" ]; then
  echo "You must specify DOCKER_IMAGE_FULL_NAME"
  exit 1
fi

echo "Deploying container image $DOCKER_IMAGE_FULL_NAME"

AZ_DEV_USER="$AZ_DEV_USER_RW"
AZ_DEV_PASSWORD="$AZ_DEV_PASSWORD_RW"

# @TODO get these values
AZ_AD_TENANT_ID=f99ef0be-7868-495e-b90c-12dee38c1fdc
AZ_RESOURCE_GROUP_NAME=rg-artnet-auto-deploy
AZ_ACI_CONTAINER_NAME=artnet-pilot-ui-playground
AZ_DNS_ZONE_NAME="az.artnet-dev.com"
AZ_DNS_ZONE_RESOURCE_GROUP=azure-dev-global-rg
DOCKER_REPOSITORY_PREFIX=artnetdev
DOCKER_REPOSITORY="$DOCKER_REPOSITORY_PREFIX.azurecr.io"

az login \
  --service-principal \
  --username "$AZ_DEV_USER" \
  --password "$AZ_DEV_PASSWORD" \
  --tenant "$AZ_AD_TENANT_ID"

echo "Deleting the previously deployed container"
az container delete \
  --resource-group "$AZ_RESOURCE_GROUP_NAME" \
  --name "$AZ_ACI_CONTAINER_NAME" \
  --yes

echo "Re-creating the container from the new image"
az container create \
  --resource-group "$AZ_RESOURCE_GROUP_NAME" \
  --name "$AZ_ACI_CONTAINER_NAME" \
  --image "$DOCKER_IMAGE_FULL_NAME" \
  --cpu 1 \
  --memory 1 \
  --registry-login-server "$DOCKER_REPOSITORY" \
  --registry-username "$DOCKER_USER_RW" \
  --registry-password "$DOCKER_PASSWORD_RW" \
  --dns-name-label "$AZ_ACI_CONTAINER_NAME" \
  --ports 80

echo "Deployed GoBot version number: $APP_VERSION"
