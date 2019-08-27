#!/bin/sh

if [ -z "$ALIBABACLOUD_ACCESS_KEY_ID" ]; then
    echo "error: ALIBABACLOUD_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$ALIBABACLOUD_SECRET_ACCESS_KEY" ]; then
    echo "error: ALIBABACLOUD_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

secrethash=$(cat $script_dir/bootstrap-master.sh | \
  sed "s/  accessKeyID: FILLIN/  accessKeyID: $(echo -n $ALIBABACLOUD_ACCESS_KEY_ID | base64)/" | \
  sed "s/  accessKeySecret: FILLIN/  accessKeySecret: $(echo -n $ALIBABACLOUD_SECRET_ACCESS_KEY | base64)/" | base64 )

cat <<EOF > $script_dir/bootstrap-master.yaml
apiVersion: v1
kind: Secret
metadata:
  name: master-user-data-secret
  namespace: default
type: Opaque
data:
  userData: |
    $secrethash
EOF
