#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  yamllint --config-data "{extends: default, rules: {line-length: {level: warning, max: 120}}}" ./examples/
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/workdir:z" \
    --entrypoint sh \
    quay.io/coreos/yamllint \
    ./hack/yaml-lint.sh
fi;
