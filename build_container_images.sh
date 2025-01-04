#!/usr/bin/bash

set -euo pipefail

main() {
    podman build -t nanokvm-backend:latest ./server

    podman build -t nanokvm-frontend:latest ./web
}

main "$@"
