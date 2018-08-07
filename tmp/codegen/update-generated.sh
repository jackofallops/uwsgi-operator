#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/sjones-sot/uwsgi-operator/pkg/generated \
github.com/sjones-sot/uwsgi-operator/pkg/apis \
sourceoftruth:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
