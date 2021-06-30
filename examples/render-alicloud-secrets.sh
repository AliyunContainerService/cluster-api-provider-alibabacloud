#!/bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "usage: $0 <filename>"
    exit 1
fi

if [ -z "$ALIBABACLOUD_ACCESS_KEY_ID" ]; then
    echo "error: ALIBABACLOUD_ACCESS_KEY_ID is not set in the environment" 2>&1
    exit 1
fi

if [ -z "$ALIBABACLOUD_SECRET_ACCESS_KEY" ]; then
    echo "error: ALIBABACLOUD_SECRET_ACCESS_KEY is not set in the environment" 2>&1
    exit 1
fi

x=$(echo -n "$ALIBABACLOUD_ACCESS_KEY_ID" | base64)
y=$(echo -n "$ALIBABACLOUD_SECRET_ACCESS_KEY" | base64)

sed -e "s/accessKeyID:.*/accessKeyID: $x/" \
    -e "s/accessKeySecret:.*/accessKeySecret: $y/" \
    "$1"
