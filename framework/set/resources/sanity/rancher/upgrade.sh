#!/bin/bash

RANCHER_CHART_REPO=$1
REPO=$2
CERT_TYPE=$3
HOSTNAME=$4
RANCHER_TAG_VERSION=$5
CHART_VERSION=$6
RANCHER_IMAGE=$7
UPGRADED_TURTLES=${8:-""}
RANCHER_AGENT_IMAGE=${9:-""}

set -ex

echo "Adding Helm chart repo"
helm repo add upgraded-rancher-${REPO} ${RANCHER_CHART_REPO}${REPO}

echo "Upgrading Rancher"
if [[ "$TURTLES" == "false" || "$TURTLES" == "toggledOn" ]]; then
    if [ "$CERT_TYPE" == "self-signed" ]; then
        if [ -n "$RANCHER_AGENT_IMAGE" ]; then
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set 'extraEnv[0].name=CATTLE_AGENT_IMAGE' \
                                                                                        --set "extraEnv[0].value=${RANCHER_AGENT_IMAGE}:${RANCHER_TAG_VERSION}" \
                                                                                        --set 'extraEnv[1].name=RANCHER_VERSION_TYPE' \
                                                                                        --set 'extraEnv[1].value=prime' \
                                                                                        --set 'extraEnv[2].name=CATTLE_BASE_UI_BRAND' \
                                                                                        --set 'extraEnv[2].value=suse' \
                                                                                        --set 'extraEnv[3].name=CATTLE_FEATURES' \
                                                                                        --set 'extraEnv[3].value=turtles=false\,embedded-cluster-api=true' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel

        else
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set 'extraEnv[0].name=CATTLE_FEATURES' \
                                                                                        --set 'extraEnv[0].value=turtles=false\,embedded-cluster-api=true' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        fi
    elif [ "$CERT_TYPE" == "lets-encrypt" ]; then
        if [ -n "$RANCHER_AGENT_IMAGE" ]; then
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set ingress.tls.source=letsEncrypt \
                                                                                        --set letsEncrypt.email=${LETS_ENCRYPT_EMAIL} \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set 'extraEnv[0].name=CATTLE_AGENT_IMAGE' \
                                                                                        --set "extraEnv[0].value=${RANCHER_AGENT_IMAGE}:${RANCHER_TAG_VERSION}" \
                                                                                        --set 'extraEnv[1].name=RANCHER_VERSION_TYPE' \
                                                                                        --set 'extraEnv[1].value=prime' \
                                                                                        --set 'extraEnv[2].name=CATTLE_BASE_UI_BRAND' \
                                                                                        --set 'extraEnv[2].value=suse' \
                                                                                        --set 'extraEnv[3].name=CATTLE_FEATURES' \
                                                                                        --set 'extraEnv[3].value=turtles=false\,embedded-cluster-api=true' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        else
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set ingress.tls.source=letsEncrypt \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set letsEncrypt.email=${LETS_ENCRYPT_EMAIL} \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set 'extraEnv[0].name=CATTLE_FEATURES' \
                                                                                        --set 'extraEnv[0].value=turtles=false\,embedded-cluster-api=true' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        fi
    else
        echo "Unsupported CERT_TYPE: $CERT_TYPE"
        exit 1
    fi
else
    if [ "$CERT_TYPE" == "self-signed" ]; then
        if [ -n "$RANCHER_AGENT_IMAGE" ]; then
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set 'extraEnv[0].name=CATTLE_AGENT_IMAGE' \
                                                                                        --set "extraEnv[0].value=${RANCHER_AGENT_IMAGE}:${RANCHER_TAG_VERSION}" \
                                                                                        --set 'extraEnv[1].name=RANCHER_VERSION_TYPE' \
                                                                                        --set 'extraEnv[1].value=prime' \
                                                                                        --set 'extraEnv[2].name=CATTLE_BASE_UI_BRAND' \
                                                                                        --set 'extraEnv[2].value=suse' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel

        else
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        fi
    elif [ "$CERT_TYPE" == "lets-encrypt" ]; then
        if [ -n "$RANCHER_AGENT_IMAGE" ]; then
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set ingress.tls.source=letsEncrypt \
                                                                                        --set letsEncrypt.email=${LETS_ENCRYPT_EMAIL} \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set 'extraEnv[0].name=CATTLE_AGENT_IMAGE' \
                                                                                        --set "extraEnv[0].value=${RANCHER_AGENT_IMAGE}:${RANCHER_TAG_VERSION}" \
                                                                                        --set 'extraEnv[1].name=RANCHER_VERSION_TYPE' \
                                                                                        --set 'extraEnv[1].value=prime' \
                                                                                        --set 'extraEnv[2].name=CATTLE_BASE_UI_BRAND' \
                                                                                        --set 'extraEnv[2].value=suse' \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        else
            helm upgrade --install rancher upgraded-rancher-${REPO}/rancher --namespace cattle-system --set global.cattle.psp.enabled=false \
                                                                                        --version ${CHART_VERSION} \
                                                                                        --set hostname=${HOSTNAME} \
                                                                                        --set rancherImage=${RANCHER_IMAGE} \
                                                                                        --set rancherImageTag=${RANCHER_TAG_VERSION} \
                                                                                        --set ingress.tls.source=letsEncrypt \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set letsEncrypt.email=${LETS_ENCRYPT_EMAIL} \
                                                                                        --set letsEncrypt.ingress.class=nginx \
                                                                                        --set agentTLSMode=system-store \
                                                                                        --devel
        fi
    else
        echo "Unsupported CERT_TYPE: $CERT_TYPE"
        exit 1
    fi
fi

echo "Waiting for Rancher to be rolled out"
kubectl -n cattle-system rollout status deploy/rancher
kubectl -n cattle-system get deploy rancher

echo "Waiting 15 seconds to be able to login to Rancher"
sleep 15