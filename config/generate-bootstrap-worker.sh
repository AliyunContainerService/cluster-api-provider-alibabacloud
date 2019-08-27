#!/bin/sh

if [ -z "$MASTER_IP" ]; then
    echo "error: MASTER_IP is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$KUBEADM_TOKEN" ]; then
    echo "error: KUBEADM_TOKEN is not set in the environment" 2>&1
    exit 1
fi


script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

secrethash=$(cat $script_dir/bootstrap-worker.sh | \
  sed "s/FILLIN:6443/$(echo  $MASTER_IP):6443/" | \
  sed "s/--token FILLIN/--token $(echo  $KUBEADM_TOKEN)/" | base64  )


cat <<EOF > $script_dir/bootstrap-worker.yaml
apiVersion: v1
kind: Secret
metadata:
  name: worker-user-data-secret
  namespace: default
type: Opaque
data:
  userData: |
    $secrethash
EOF
