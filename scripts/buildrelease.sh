#!/usr/bin/env bash

set -euo pipefail

APP_NAME="leafy"
DIST_DIR="dist"
if [[ "${GITHUB_REF_TYPE:-}" == "tag" ]]; then
  VERSION="${GITHUB_REF_NAME}"
else
  VERSION="$(git describe --tags --exact-match 2>/dev/null || echo "dev")"
fi

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

build() {
  local goarch="$1"
  local output="${APP_NAME}-linux-${goarch}"

  echo "Building ${output} with version ${VERSION}..."

  GOOS="linux" GOARCH="${goarch}" CGO_ENABLED=0 \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o "${DIST_DIR}/${output}" .
}

build amd64
build arm64

echo "Generating checksums..."

(
  cd "${DIST_DIR}"
  sha256sum * > checksums.txt
)

echo "Linux release artifacts created in ${DIST_DIR}/"