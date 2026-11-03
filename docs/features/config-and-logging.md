# Configuration & Logging

> 🚧 **Status: not wired in.** These packages build and pass checks, but `main.go` and `app/tui` never call them.
> Leafy currently writes no log file and reads no configuration.

## Configuration (`internal/config`)

The package doc says it is meant to "manage a user's configs for ease-of-use and consistency in multiple sessions".

```go
type Config struct {
    LogPath string
}

func Default() *Config // LogPath = "$HOME/.local/state/leafy/logs/leafy.log"
```

The default path follows the XDG Base Directory convention for state data (`~/.local/state`). The config is only
built in code; nothing loads it from a file.

## Logging (`internal/log`)

```go
func InitializeLogger(cfg *config.Config) (*slog.Logger, CloseFunc, error)
```

- Creates the log directory with mode `0700`.
- Opens `cfg.LogPath` with `O_APPEND|O_CREATE|O_WRONLY` and mode `0600`.
- Wraps the file in an 8 KiB `bufio.Writer`.
- Returns a `log/slog` logger with a **JSON** handler at level `Info`, plus a `CloseFunc` that closes the file.

### Example (intended usage)

```go
cfg := config.Default()
logger, closeLog, err := log.InitializeLogger(cfg)
if err != nil { /* handle */ }
defer closeLog()
logger.Info("device mounted", "path", "/dev/sdb1", "mountpoint", "/run/media/me/CAM")
```

> ⚠️ Before wiring this in, note that `CloseFunc` closes the file **without flushing** the `bufio.Writer`. Up to
> 8 KiB of buffered log records can be lost on exit. Flush the buffer before closing the file.
