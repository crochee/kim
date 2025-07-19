#!/bin/bash

# This shell is used to auto generate some useful tools for k8s, such as clientset, lister, informer and so on.
# We don't use this tool to generate deepcopy because kubebuilder (controller-tools) has coverred that part.

set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "${0}")/.."

# kube_codegen.sh may be confused by the symbolic links in the path, so cd to the canonicalized path
REPO_ROOT="$(readlink -f .)"
cd "${REPO_ROOT}"

CODEGEN_VERSION=$(grep 'k8s.io/code-generator' go.sum | awk '{print $2}' | sed 's/\/go.mod//g' | head -1)
CODEGEN_PKG=$(echo $(go env GOPATH)"/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}")
if [[ ! -d ${CODEGEN_PKG} ]]; then
    echo "${CODEGEN_PKG} is missing. Running 'go mod download'."
    go mod download
fi
echo ">> Using ${CODEGEN_PKG}"

# 获取go mod的包名
THIS_PKG="$(go list -m | cut -d' ' -f3)"

# shellcheck source=/dev/null
source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_helpers "${REPO_ROOT}/api" \
    --boilerplate "${REPO_ROOT}/hack/boilerplate.go.txt"

kube::codegen::gen_client "${REPO_ROOT}/api" \
    --with-watch \
    --with-applyconfig \
    --output-dir "${REPO_ROOT}/pkg/client" \
    --boilerplate "${REPO_ROOT}/hack/boilerplate.go.txt" \
    --output-pkg "${THIS_PKG}/pkg/client"
