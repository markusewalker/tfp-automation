#!/bin/bash

set -x
set -eu

DEBUG="${DEBUG:-false}"
TERRAFORM_VERSION="${TERRAFORM_VERSION:-}"
RKE_PROVIDER_VERSION="${RKE_PROVIDER_VERSION:-}"
RANCHER2_PROVIDER_VERSION="${RANCHER2_PROVIDER_VERSION:-}"
LOCALS_PROVIDER_VERSION="${LOCALS_PROVIDER_VERSION:-}"
CLOUD_PROVIDER_VERSION="${CLOUD_PROVIDER_VERSION:-}"
KUBERNETES_PROVIDER_VERSION="${KUBERNETES_PROVIDER_VERSION:-}"
LETS_ENCRYPT_EMAIL="${LETS_ENCRYPT_EMAIL:-}"

TRIM_JOB_NAME=$(basename "$JOB_NAME")

if [ "false" != "${DEBUG}" ]; then
    echo "Environment:"
    env | sort
fi

count=0
while [[ 3 -gt $count ]]; do
    docker build . -f Dockerfile --build-arg CONFIG_FILE=config.yml --build-arg PEM_FILE=key.pem \
                                                                    --build-arg TERRAFORM_VERSION="$TERRAFORM_VERSION" \
                                                                    --build-arg RKE_PROVIDER_VERSION="$RKE_PROVIDER_VERSION" \
                                                                    --build-arg RANCHER2_PROVIDER_VERSION="$RANCHER2_PROVIDER_VERSION" \
                                                                    --build-arg LOCALS_PROVIDER_VERSION="$LOCALS_PROVIDER_VERSION" \
                                                                    --build-arg CLOUD_PROVIDER_VERSION="$CLOUD_PROVIDER_VERSION" \
                                                                    --build-arg KUBERNETES_PROVIDER_VERSION="$KUBERNETES_PROVIDER_VERSION" \
                                                                    --build-arg LETS_ENCRYPT_EMAIL="$LETS_ENCRYPT_EMAIL" \
                                                                    -t tfp-automation-validation-"${TRIM_JOB_NAME}""${BUILD_NUMBER}"

    if [[ $? -eq 0 ]]; then break; fi
    count=$(($count + 1))
    echo "Repeating failed Docker build ${count} of 3..."
done