#!/bin/bash

PEM_FILE=$1
USER=$2
GROUP=$3
NODE_PRIVATE_IP=$4
REGISTRATION_COMMAND=$5
REGISTRY=$6

echo ${PEM_FILE} | sudo base64 -d > /home/${USER}/airgap.pem
echo "${REGISTRATION_COMMAND}" > /home/${USER}/registration_command.txt
REGISTRATION_COMMAND=$(cat /home/$USER/registration_command.txt)

PEM=/home/${USER}/airgap.pem
sudo chmod 600 ${PEM}
sudo chown ${USER}:${GROUP} ${PEM}

runSSH() {
  local server="$1"
  local cmd="$2"

  ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=300 -o ServerAliveInterval=30 -o ServerAliveCountMax=10 -i "$PEM" "$USER@$server" \
  "export USER=${USER}; \
   export GROUP=${GROUP}; \
   export NODE_PRIVATE_IP=${NODE_PRIVATE_IP}; \
   export REGISTRY=${REGISTRY}; \
   export REGISTRATION_COMMAND=${REGISTRATION_COMMAND}; \
   export REGISTRY=${REGISTRY}; $cmd"
}

setupDockerDaemon() {
  echo "{ \"insecure-registries\" : [ \"${REGISTRY}\" ] }" | sudo tee /etc/docker/daemon.json > /dev/null
}

dockerDaemonFunction=$(declare -f setupDockerDaemon)
runSSH "${NODE_PRIVATE_IP}" "${dockerDaemonFunction}; setupDockerDaemon"

runSSH "${NODE_PRIVATE_IP}" "sudo systemctl daemon-reload && sudo systemctl restart docker"

MAX_RETRIES=5
RETRY_DELAY=15
ATTEMPT=1
SUCCESS=0

while [ $ATTEMPT -le $MAX_RETRIES ]; do
  runSSH "${NODE_PRIVATE_IP}" "${REGISTRATION_COMMAND}"
  EXIT_CODE=$?

  if [ $EXIT_CODE -eq 0 ]; then
    SUCCESS=1
    break
  else
    sleep $RETRY_DELAY
    ATTEMPT=$((ATTEMPT+1))
  fi
done

if [ $SUCCESS -ne 1 ]; then
  exit 1
fi