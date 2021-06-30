#!/bin/sh

if [ -z "$alibabacloud_ACCESS_KEY_ID" ]; then
    echo "error: alibabacloud_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$alibabacloud_SECRET_ACCESS_KEY" ]; then
    echo "error: alibabacloud_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

secrethash=$(cat $script_dir/bootstrap.sh | \
  sed "s/  accessKeyID: FILLIN/  accessKeyID: $(echo -n $alibabacloud_ACCESS_KEY_ID | base64)/" | \
  sed "s/  accessKeySecret: FILLIN/  accessKeySecret: $(echo -n $alibabacloud_SECRET_ACCESS_KEY | base64)/" )

cat <<EOF > $script_dir/bootstrap.yaml
apiVersion: v1
kind: Secret
metadata:
  name: alibabacoud-credentials-secret
  namespace: default
type: Opaque
data:
  userData: $secrethash
EOF
