# Features

| Feature                                        | Status        | Code                                               |
|------------------------------------------------|---------------|----------------------------------------------------|
| [Device Management](device-management.md)      | ✅ Active      | `internal/device`, `app/tui/cmd.go`                |
| [Media Discovery](media-discovery.md)          | ✅ Active      | `internal/file/file.go`, `app/tui/format.go`       |
| [File Transfer](file-transfer.md)              | ✅ Active      | `internal/file/copy.go`, `internal/file/progress.go` |
| [Terminal UI](terminal-ui.md)                  | ✅ Active      | `app/tui`, `internal/style`                        |
| [Versioning](versioning.md)                    | ✅ Active      | `main.go`, `scripts/buildrelease.sh`               |
| [Configuration & Logging](config-and-logging.md) | 🚧 Not wired in | `internal/config`, `internal/log`               |
| Device power-off                               | 🚧 Not wired in | `device.PowerOffDevice` (see [Device Management](device-management.md#power-off-not-yet-used)) |
| Loading spinner / minimum durations            | 🚧 Not wired in | `internal/style/loading.go`, `style.LoadDelay`, `style.QuitDelay` |
