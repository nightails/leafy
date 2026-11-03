# Known Limitations & Roadmap

This page lists behaviour you may trip over, found by reading the current source. Each entry names the code
responsible, so it can be fixed.

## Bugs / sharp edges

| # | Area | Description | Code |
|---|------|-------------|------|
| 1 | Input | The `space` key is handled globally as "toggle selection", whatever the step. In the **destination** step you cannot type spaces into the path, and pressing space silently toggles a list item instead. | `app/tui/tui.go`, `case " "` |
| 2 | Input | Pressing `space` while the media list is **empty** (no USB device or no media found) panics. `SelectedItem()` returns `nil`, and `item.(mediumItem)` fails on it. Bubble Tea restores the terminal and prints the panic. | `app/tui/tui.go`, `case " "` |
| 3 | Exit | When a copy fails, `copyErrorMsg` sets `m.err` but leaves `currStep == copying`. `ctrl+c` is ignored in that step, so **the app cannot be quit normally**. You have to kill the process and unmount manually. | `app/tui/tui.go`, `case "ctrl+c"` / `copyErrorMsg` |
| 4 | Copy | On copy or sync failure, `os.Remove(tmp)` is called with the bare temp **file name**, so it resolves against the process working directory instead of the destination directory. The `<name>.tmp` file stays behind, and the open handle `out` is never closed on those paths. | `internal/file/copy.go` |
| 5 | Unmount | On exit, **every** detected USB partition is unmounted, including ones that were mounted before Leafy started. The README says only devices "mounted by the app" are unmounted. | `app/tui/cmd.go`, `initDevicesCmd` / `removeDevicesCmd` |
| 6 | Unmount | `removeDevicesCmd` is batched with `tea.Quit`, so unmount errors are never shown. The loop also stops at the first failure. | `app/tui/tui.go`, `app/tui/cmd.go` |
| 7 | Copy | The destination is flat (`dest/<name>`), and existing files are **overwritten**. Files with the same name from different folders or devices clobber each other. | `app/tui/tui.go` (`TODO: apply naming scheme with increment`) |
| 8 | Mount | If mounting any single partition fails, startup aborts for all devices. | `app/tui/cmd.go`, `initDevicesCmd` |
| 9 | UX | No loading indicator during detection and scanning, and no "no devices / no media found" message. An empty list is all you see. | `app/tui/tui.go`, `devicesMsg` / `mediaMsg` |
| 10 | UX | You cannot go back a step, cancel a running copy, or change the destination after confirming it. | `app/tui/tui.go` (`TODO: implement canceling`) |
| 11 | Scan | Symlinks and macOS `._*` sidecar files with media extensions are listed. Symlinks then fail to copy, because only regular files are accepted. | `internal/file/file.go` |
| 12 | Logging | `CloseFunc` closes the log file without flushing its `bufio.Writer`, so buffered records are lost. The logger is not used yet. | `internal/log/log.go` |
| 13 | View | The `Selected media has been deleted.` message only shows when `finished` is reached with zero media, which the current flow cannot produce. | `app/tui/tui.go`, `View` |
| 14 | Release | `cd.yml` releases for any successful CI run whose `head_branch` starts with `v`. That includes branches, not only tags. | `.github/workflows/cd.yml` |

## Platform limitations

- Linux only: it depends on `lsblk`, `udisksctl`/UDisks2 and `golang.org/x/sys/unix`.
- Only devices with `TRAN=usb` are detected. Internal SD card readers that sit on the PCI or MMC bus (for example
  `/dev/mmcblk0`) are **not** picked up.
- Mounting needs polkit authorisation for the current user, which is normally available only in a local desktop
  session.
- Destination directories are created with mode `0700`, so only the owner can access them.

## Code that exists but is not wired in yet

| Item | Location | Probable purpose |
|------|----------|------------------|
| `config.Config`, `config.Default()` | `internal/config` | Persisted user settings, starting with the log path |
| `log.InitializeLogger` | `internal/log` | JSON log file at `~/.local/state/leafy/logs/leafy.log` |
| `device.PowerOffDevice` | `internal/device/device.go` | Safely power off the device after transfer |
| `style.NewLineSpinner`, `style.MinDuration`, `style.LoadDelay`, `style.QuitDelay` | `internal/style` | Loading and quitting screens with a minimum display time |
| `deleteDoneMsg` | `app/tui/msg.go` | Optional "delete source media after copy" step |
| `USBDevice.Label` | `internal/device/device.go` | Showing partition labels in the UI |
| `assets/leafy_ascii_art.txt`, `assets/tui_layout.txt` | `assets/` | Splash logo and layout design (the layout file is empty) |

## Testing

There are no automated tests (`*_test.go`). CI runs `go test -cover ./...`, which passes with 0% coverage. See the
[Development Guide](development.md#quality-checks-same-as-ci) for suggested first tests.

## Suggested roadmap (derived from TODOs and unused code)

1. Fix the input handling: handle `space` only in the media step and guard against an empty list (bugs 1–2).
2. Allow quitting after a copy error, and clean up temp files correctly (bugs 3–4).
3. Track which devices Leafy mounted itself and unmount only those. Optionally power them off (bugs 5–6).
4. Name collision handling with incrementing suffixes, and optionally keep the folder structure (bug 7).
5. Copy cancellation, a loading spinner, and empty-state messages (bugs 9–10).
6. Wire in config and logging, with a flush fix (bug 12).
7. Add unit tests for `internal/file`, `internal/device` and the TUI state machine.
