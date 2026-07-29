#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'Required command not found: %s\n' "$1" >&2
        exit 1
    fi
}

require_command ansible-playbook
require_command docker
require_command go
require_command terraform

mapfile -t go_files < <(find "$ROOT_DIR/platform-api" -name '*.go' -type f -print)
unformatted_files="$(gofmt -l "${go_files[@]}")"

if [[ -n "$unformatted_files" ]]; then
    printf 'Go files need formatting:\n%s\n' "$unformatted_files" >&2
    exit 1
fi

(
    cd "$ROOT_DIR/platform-api"
    go test ./...
)

bash -n "$ROOT_DIR"/platform-lab/scripts/*.sh

ansible-playbook \
    -i "$ROOT_DIR/platform-lab/ansible/inventory.ini" \
    "$ROOT_DIR/platform-lab/ansible/playbook.yml" \
    --syntax-check

terraform -chdir="$ROOT_DIR/platform-lab/terraform" fmt -check -recursive
TF_IN_AUTOMATION=1 terraform -chdir="$ROOT_DIR/platform-lab/terraform" init -backend=false -input=false
terraform -chdir="$ROOT_DIR/platform-lab/terraform" validate

docker build --tag platform-api:ci "$ROOT_DIR/platform-api"
