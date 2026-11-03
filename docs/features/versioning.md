# Versioning

Code: `main.go` (`version`, `getVersion`) and `scripts/buildrelease.sh`.

The header shows `Leafy <version>`. `getVersion()` picks the first value that applies:

| Priority | Source                                         | When it applies                                              | Example        |
|----------|------------------------------------------------|--------------------------------------------------------------|----------------|
| 1        | `main.version`, injected with `-ldflags -X`     | Release builds from `scripts/buildrelease.sh` / CI           | `v0.3.0`       |
| 2        | `debug.ReadBuildInfo().Main.Version`           | `go install github.com/nightails/leafy@vX.Y.Z` (or `@latest`) | `v0.3.0`       |
| 3        | Fallback `"dev"`                                | `go run`, or `go build` in a working copy (`(devel)`)         | `dev`          |

## How the release script resolves the version

```bash
if [[ "${GITHUB_REF_TYPE:-}" == "tag" ]]; then
  VERSION="${GITHUB_REF_NAME}"
else
  VERSION="$(git describe --tags --exact-match 2>/dev/null || echo "dev")"
fi
```

Releases are built by a `workflow_run`-triggered job that checks out the tag. The `git describe --tags --exact-match`
branch therefore usually supplies the tag name, such as `v1.2.3`. The script embeds it with
`-ldflags="-s -w -X main.version=${VERSION}"`.

To build a versioned binary locally:

```bash
go build -ldflags "-X main.version=v0.0.0-local" -o leafy .
```
