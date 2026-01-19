#!/usr/bin/env bash

usage() {
    local script
    script=$(basename "$0")
    echo "usage:" >&2
    echo "$script" >&2
    exit "${1:-1}"
}

check_dependent_tools() {
    local missing=()
    for tool in "${@}"; do
        if ! command -v "${tool}" &>/dev/null; then
            missing+=("$tool")
        fi
    done

    if ((${#missing[@]})); then
        echo "error:missing required tool(s):${missing[*]}" >&2
        exit 1
    fi
}

check_parameters() {
    if (("$#" != 0)); then
        usage
    fi
}

process_opts() {
    while getopts ":h" opt; do
        case $opt in
        h)
            usage 0
            ;;
        *)
            echo "error:unsupported option -$opt" >&2
            usage
            ;;
        esac
    done
}

main() {
    REQUIRED_TOOLS=()
    check_dependent_tools "${REQUIRED_TOOLS[@]}"
    check_parameters "${@}"
    OPTIND=1
    process_opts "${@}"
    shift $((OPTIND - 1))

    local ROOT_DIR="${HOME}"/repos/go-tools/
    local BIN_DIR="${ROOT_DIR}"/bin
    mkdir -p "$BIN_DIR"

    echo "Building all tools in go-tools/tools/..."

    for dir in tools/*; do
        if [ -f "$dir/main.go" ]; then
            echo "Building $(basename ${dir})..."
            (cd "$dir" && go build -o "$BIN_DIR/$(basename "${dir}")" main.go)
        fi
    done
}

main "${@}"
