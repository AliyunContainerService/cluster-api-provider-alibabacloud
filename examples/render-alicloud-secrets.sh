#!/bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "usage: $0 <filename>"
    exit 1
fi

if [ -z "$ALICLOUD_ACCESS_KEY_ID" ]; then
    echo "error: ALICLOUD_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$ALICLOUD_SECRET_ACCESS_KEY" ]; then
    echo "error: ALICLOUD_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

x=$(echo -n "$ALICLOUD_ACCESS_KEY_ID" | base64)
y=$(echo -n "$ALICLOUD_SECRET_ACCESS_KEY" | base64)

sed -e "s/alicloud_access_key_id:.*/alicloud_access_key_id: $x/" \
    -e "s/alicloud_secret_access_key:.*/alicloud_secret_access_key: $y/" \
    "$1"
