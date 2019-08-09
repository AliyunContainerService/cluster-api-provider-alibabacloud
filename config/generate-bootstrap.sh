#!/bin/sh

if [ -z "$ALICLOUD_ACCESS_KEY_ID" ]; then
    echo "error: ALICLOUD_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$ALICOUD_SECRET_ACCESS_KEY" ]; then
    echo "error: ALICLOUD_ACCESS_KEY_SECRET is not set in the environment" 2>&1
    exit 1
fi

script_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)"

secrethash=$(cat $script_dir/bootstrap.sh | \
  sed "s/  alicloud_access_key_id: FILLIN/  alicloud_access_key_id: $(echo -n $ALICLOUD_ACCESS_KEY_ID | base64)/" | \
  sed "s/  ALICLOUD_ACCESS_KEY_SECRET: FILLIN/  ALICLOUD_ACCESS_KEY_SECRET: $(echo -n $ALICLOUD_ACCESS_KEY_SECRET | base64)/" | \
  base64 --w=0)

cat <<EOF > $script_dir/bootstrap.yaml
apiVersion: v1
kind: Secret
metadata:
  name: master-user-data-secret
  namespace: default
type: Opaque
data:
  userData: $secrethash
EOF
