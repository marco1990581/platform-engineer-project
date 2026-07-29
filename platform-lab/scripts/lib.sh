#!/usr/bin/env bash


set -euo pipefail

readonly GREEN="\033[0;32m"
readonly YELLOW="\033[1,33m"
readonly RED="\033[0,31m"
readonly BLUE="\033[0,34m"
readonly NC="\033[0m"

info(){
    printf "${BLUE}[INFO]${NC} %s\n" "$*"
}

success(){
    printf "${GREEN}[ OK ]${NC} %s\n" "$*"
}

warn(){
    printf "${YELLOW}[ WARN ]${NC} %s\n" "$*"
}

error(){
     printf "${RED}[ FAIL ]${NC} %s\n" "$*" >&2
}


command_exists() {
    command -v "$1" >/dev/null 2>&1
}

require_apt_host() {
    if [[ ! -r /etc/os-release ]]; then
        error "This installer requires a Debian or Ubuntu host."
        exit 1
    fi

    . /etc/os-release

    case "$ID" in
        debian | ubuntu)
            ;;
        *)
            error "Unsupported distribution: $ID"
            exit 1
            ;;
    esac
}

require_command() {
    if ! command_exists "$1"; then
        error "Required command not found: $1"
        exit 1
    fi
}

linux_arch() {
    case "$(uname -m)" in
        x86_64)
            printf 'amd64\n'
            ;;
        aarch64 | arm64)
            printf 'arm64\n'
            ;;
        *)
            error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

require_sudo(){

   sudo -v
}

install_package() {
    local package
    local missing_packages=()

    for package in "$@"; do
        if dpkg -s "$package" >/dev/null 2>&1; then
            info "$package already installed."
        else
            missing_packages+=("$package")
        fi
    done

    if ((${#missing_packages[@]} == 0)); then
        return
    fi

    info "Installing: ${missing_packages[*]}"
    sudo apt install -y "${missing_packages[@]}"
    success "Package installation completed."
}
