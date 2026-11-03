# Getting Started

## Requirements

Leafy runs on **Linux only**. It calls Linux-specific tools and syscalls (`golang.org/x/sys/unix`), and it uses
UDisks2 for mounting.

| Requirement      | Why it is needed                                                                    |
|------------------|-------------------------------------------------------------------------------------|
| Linux (amd64 / arm64) | Release binaries are published for these two architectures                     |
| `lsblk` (util-linux)  | Detects USB disks and partitions (must support `-J` JSON output and `MOUNTPOINTS`) |
| `udisksctl` (udisks2) | Mounts and unmounts partitions as a regular user                               |
| A running UDisks2 daemon / polkit session | Lets `udisksctl mount` work without `sudo`             |
| Go 1.25+ (only for `go install` / building from source) | `go.mod` declares `go 1.25`. The code uses `os.Root` APIs added in recent Go versions |

Most desktop distributions (Fedora, Ubuntu, Arch with a desktop environment, and others) ship `lsblk` and
`udisksctl` out of the box. On minimal systems, install `util-linux` and `udisks2`.

## Installation

### Option 1: `go install`

```bash
go install github.com/nightails/leafy@latest
```

The binary lands in `$(go env GOPATH)/bin`. Make sure that directory is on your `PATH`. A binary installed this way
reports the module version (for example `v0.2.0`) in its header. See [Versioning](features/versioning.md).

### Option 2: Pre-built binary

1. Download `leafy-linux-amd64` or `leafy-linux-arm64` from the GitHub **Releases** page.
2. Optionally verify it against `checksums.txt`:
   ```bash
   sha256sum -c checksums.txt --ignore-missing
   ```
3. Make it executable and run it:
   ```bash
   chmod +x leafy-linux-amd64
   ./leafy-linux-amd64
   ```

### Option 3: Build from source

```bash
git clone https://github.com/nightails/leafy.git
cd leafy
go build -o leafy .
./leafy
```

See the [Development Guide](development.md) for more.

## First run

1. Plug in your camera, microphone, or other USB storage device.
2. Run `leafy` in a terminal.
3. Leafy finds and mounts the device, then lists the media files it found.
4. Select files, type a destination folder, and press `enter` to copy.
5. When copying ends, Leafy unmounts the device and quits on its own.

The full walkthrough is in the [User Workflow](user-workflow.md) guide.

> **Note:** Leafy takes no command-line flags or arguments. It always starts the interactive TUI.
