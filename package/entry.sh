#!/bin/bash
set -ex

MD=${RANCHER_METADATA_ADDRESS:-169.254.169.250}
while ! curl -s -f http://${MD}/2015-12-19/stacks/Kubernetes/services/kubernetes/uuid; do
    echo Waiting for metadata
    sleep 1
done

/usr/bin/update-rancher-ssl

CATTLE_ACCESS_KEY=$CATTLE_ENVIRONMENT_ADMIN_ACCESS_KEY
CATTLE_SECRET_KEY=$CATTLE_ENVIRONMENT_ADMIN_SECRET_KEY

UUID=$(curl -s http://rancher-metadata/2015-12-19/stacks/Kubernetes/services/kubernetes/uuid)
ACTION=$(curl -s -u $CATTLE_ACCESS_KEY:$CATTLE_SECRET_KEY "$CATTLE_URL/services?uuid=$UUID" | jq -r '.data[0].actions.certificate')
KUBERNETES_URL=${KUBERNETES_URL:-https://kubernetes:6443}

if [ -n "$ACTION" ]; then
    mkdir -p /etc/kubernetes/ssl
    cd /etc/kubernetes/ssl
    curl -s -u $CATTLE_ACCESS_KEY:$CATTLE_SECRET_KEY -X POST $ACTION > certs.zip
    unzip -o certs.zip

    export CATTLE_ACCESS_KEY=$CATTLE_AGENT_ACCESS_KEY
    export CATTLE_SECRET_KEY=$CATTLE_AGENT_SECRET_KEY

    TOKEN=$(cat /etc/kubernetes/ssl/key.pem | sha256sum | awk '{print $1}')

    cat > /etc/kubernetes/ssl/kubeconfig << EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    certificate-authority: /etc/kubernetes/ssl/ca.pem
    server: "$KUBERNETES_URL"
  name: "Default"
contexts:
- context:
    cluster: "Default"
    user: "Default"
  name: "Default"
current-context: "Default"
users:
- name: "Default"
  user:
    token: "$TOKEN"
EOF

    exec "$@"
fi
