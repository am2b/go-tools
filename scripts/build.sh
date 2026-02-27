#!/usr/bin/env bash

usage() {
    local script
    script=$(basename "$0")
    echo "编译指定的工具"
    echo "usage:" >&2
    echo "$script tool_name" >&2
    exit "${1:-1}"
}

check_dependent_tools() {
    local missing=()
    for tool in "${@}"; do
        if ! command -v "${tool}" &> /dev/null; then
            missing+=("$tool")
        fi
    done

    if ((${#missing[@]})); then
        echo "error:missing required tool(s):${missing[*]}" >&2
        exit 1
    fi
}

check_parameters() {
    if (("$#" != 1)); then
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

    local tool="${1}"

    local ROOT_DIR="${HOME}"/repos/go-tools/
    local BIN_DIR="${ROOT_DIR}"/bin
    mkdir -p "$BIN_DIR"

    echo "Building ${tool}..."

    if [ -f "./tools/${tool}/main.go" ]; then
        (cd "./tools/${tool}" && go build -o "$BIN_DIR/tool" main.go)
    fi

    echo "Done"
}

main "${@}"
