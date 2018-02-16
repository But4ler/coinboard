#!/bin/bash -xe

CWD="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# unset the docker host variable to avoid building on clusters
unset DOCKER_HOST

function docker_tag {
  IMAGE_ID=$1
  IMAGE_NAME=$2
  docker tag ${IMAGE_ID} arnobroekhof/${IMAGE_NAME}:armhf
}

function docker_build {
  TAG=$1
  DOCKER_FILE=$2
  BUILD_PATH=$3
  IMAGE_NAME=$4
  docker build --force-rm -t ${TAG} -f ${DOCKER_FILE} ${BUILD_PATH}
  docker_tag ${TAG} ${IMAGE_NAME}
  docker_push ${TAG} ${IMAGE_NAME}
}

function docker_push {
  IMAGE_ID=$1
  IMAGE_NAME=$2
  docker push arnobroekhof/${IMAGE_NAME}:armhf
}

CUR_DATE=`date +%Y%m%d%H%M`
docker_build prometheus prometheus/Dockerfile.prometheus-armhf ${CWD}/prometheus/ prometheus
docker_build alertmanager prometheus/Dockerfile.alertmanager-armhf ${CWD}/prometheus/ alertmanager
