#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CLUSTER_NAME="platform-lab"
readonly CONFIG_FILE="$SCRIPT_DIR/../kubernetes/kind/cluster.yaml"

bash "$SCRIPT_DIR/00-preflight.sh"

if ! command -v docker >/dev/null 2>&1; then
    printf 'Docker must be installed before creating the Kind cluster.\n' >&2
    exit 1
fi

if ! docker info >/dev/null 2>&1; then
    printf 'Docker is installed but unavailable to the current user.\n' >&2
    exit 1
fi

if ! command -v kind >/dev/null 2>&1; then
    printf 'Kind must be installed before creating the cluster.\n' >&2
    exit 1
fi

if ! command -v kubectl >/dev/null 2>&1; then
    printf 'kubectl must be installed before creating the cluster.\n' >&2
    exit 1
fi

if kind get clusters | grep -Fxq "$CLUSTER_NAME"; then
    printf 'Kind cluster %q already exists.\n' "$CLUSTER_NAME"
    exit 0
fi

kind create cluster --config "$CONFIG_FILE"
kubectl cluster-info --context "kind-${CLUSTER_NAME}"
