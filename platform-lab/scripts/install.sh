#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/lib.sh"

bash "$SCRIPT_DIR/00-preflight.sh"
require_sudo

sudo apt update
install_package ansible-core

bootstrap_user="${SUDO_USER:-$USER}"
bootstrap_home="$(getent passwd "$bootstrap_user" | cut -d: -f6)"

if [[ -z "$bootstrap_home" ]]; then
    error "Unable to determine the home directory for $bootstrap_user."
    exit 1
fi

ansible-playbook \
    -i "$SCRIPT_DIR/../ansible/inventory.ini" \
    "$SCRIPT_DIR/../ansible/playbook.yml" \
    --extra-vars "bootstrap_user=$bootstrap_user bootstrap_home=$bootstrap_home"
