#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$SCRIPT_DIR/lib.sh"

require_apt_host

memory_kb="$(awk '/MemTotal:/ { print $2 }' /proc/meminfo)"
available_disk_kb="$(df --output=avail -k "$HOME" | tail -n 1)"

info "Detected $(. /etc/os-release && printf '%s %s' "$ID" "$VERSION_ID") on $(linux_arch)."

if ((memory_kb < 4194304)); then
    warn "Less than 4 GiB RAM detected. Install tools now, but create the Kind cluster only when more memory is available."
fi

if ((available_disk_kb < 10485760)); then
    warn "Less than 10 GiB free in $HOME. Free disk space before creating a Kind cluster."
fi

success "Preflight checks completed."
