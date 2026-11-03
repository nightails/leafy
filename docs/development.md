# Development Guide

## Prerequisites

| Tool          | Version / note                                                             |
|---------------|-----------------------------------------------------------------------------|
| Go            | ≥ 1.25 (`go.mod`). CI uses **1.26.0**                                      |
| `lsblk`       | From util-linux, with JSON output (`-J`) and the `MOUNTPOINTS` column        |
| `udisksctl`   | From udisks2, with a running UDisks2 daemon                                 |
| `staticcheck` | Optional locally: `go install honnef.co/go/tools/cmd/staticcheck@latest`    |
| `gosec`       | Optional locally: `go install github.com/securego/gosec/v2/cmd/gosec@latest` |

## Setup, run, build

```bash
git clone https://github.com/nightails/leafy.git
cd leafy
go mod download

# run (version shows as "dev")
go run .

# build
go build -o leafy .

# build like a release (stripped, versioned, static)
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=v0.0.0-local" -o leafy .

# build all release artifacts into dist/
bash scripts/buildrelease.sh
```

The `leafy` binary is listed in `.gitignore`. The `dist/` directory is not, so do not commit it.

## Quality checks (same as CI)

```bash
go test -cover ./...          # unit tests (+ coverage)
gosec ./...                   # security scan
test -z "$(go fmt ./...)"     # formatting must be clean
staticcheck ./...             # lint
```

> There are currently **no `_test.go` files**. `go test` succeeds, but every package reports 0% coverage. Good
> first candidates for tests:
> - `internal/file`: `GetFiles` (with `t.TempDir()`), `CopyWithProgress` (content, permissions, timestamps, root
>   escape rejection), and the progress callback.
> - `internal/device`: `FindUSBDevices` filtering, after refactoring `readLSBLK` so JSON can be injected;
>   `validateBlockDevice`; and the `udisksctl` output regexes.
> - `app/tui`: `formatFileSize`, `expandHome`, and `Update` transitions driven by synthetic messages.

## Manual testing tips

- Use a USB stick holding a few small `.mp4` / `.wav` files, some with upper-case extensions and some in nested
  folders.
- Try the device both pre-mounted (opened in your file manager first) and unmounted.
- Check copied files with `stat`: the mode and timestamps should match the source.
- A loop device does **not** work as a substitute, because detection requires `TRAN=usb`.

## Code conventions

- Format with `gofmt`. Keep `staticcheck` and `gosec` clean, and justify any `#nosec` with a comment, as in
  `udisksctl.go`.
- Do all blocking work in `tea.Cmd` functions (`app/tui/cmd.go`). Never block inside `Update`.
- Add new message types to `app/tui/msg.go` and handle them in `Model.Update`.
- Keep system interaction in `internal/*`. `app/tui` should only orchestrate.
- Wrap errors with context, for example `fmt.Errorf("open src: %w", err)`.

## Adding a feature: checklist

1. Put the domain logic in an `internal/` package.
2. Wrap it in a `tea.Cmd` that returns a dedicated message.
3. Handle the message in `Update`. If the feature needs a new screen, add a `step`.
4. Render it in `View()` and update the help line.
5. Document it in `docs/features/` and update `docs/user-workflow.md`.
