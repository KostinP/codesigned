#!/bin/bash

# Environment variables loader

set -e

ENV_FILE="${1:-.env}"

load_env() {
    if [ -f "$ENV_FILE" ]; then
        echo "Loading environment from $ENV_FILE"
        set -a
        source "$ENV_FILE"
        set +a
    else
        echo "Warning: $ENV_FILE not found"
    fi
}

get_env() {
    local var_name="$1"
    local default_value="${2:-}"
    
    if [ -f "$ENV_FILE" ]; then
        source "$ENV_FILE"
    fi
    
    local value="${!var_name:-$default_value}"
    echo "$value"
}

check_required_vars() {
    local missing_vars=()
    
    for var in "$@"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        echo "Error: Missing required environment variables: ${missing_vars[*]}"
        echo "Please check your $ENV_FILE file"
        exit 1
    fi
}

# Load environment when script is sourced
if [ "${BASH_SOURCE[0]}" = "$0" ]; then
    # Script is being executed
    case "$1" in
        "get")
            get_env "$2" "$3"
            ;;
        "check")
            shift
            load_env
            check_required_vars "$@"
            ;;
        *)
            load_env
            ;;
    esac
else
    # Script is being sourced
    load_env
fi