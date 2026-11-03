# File Transfer

Code: `internal/file/copy.go` (`CopyWithProgress`) and `internal/file/progress.go`. Driven by `copyMediaCmd` and
`listenCopyCmd` in `app/tui/cmd.go`.

## Destination resolution (TUI)

When the user confirms the destination path:

1. `expandHome` turns `~` into `$HOME` and `~/x` into `$HOME/x`. Other forms, like `~user`, are left unchanged.
2. For every selected medium, `dest = filepath.Join(<destination>, <file name>)`. The layout is **flat**: the folder
   structure from the device is dropped.
3. Files are copied **sequentially**, in list order.

## `CopyWithProgress(srcRoot, src, dst, progress)`

```mermaid
flowchart TD
    A[Lstat src] -->|not regular file| X[error]
    A --> B[Read atime/mtime via unix.Stat]
    B --> C[MkdirAll dst dir 0700]
    C --> D[Compute src path relative to srcRoot<br/>reject '.', '..', '../*']
    D --> E[os.OpenRoot srcRoot → Open rel]
    E --> F[os.OpenRoot dstDir → create name.tmp<br/>O_CREATE|O_TRUNC, src perms]
    F --> G[io.CopyBuffer 1 MiB<br/>via progressWriter]
    G --> H[fsync tmp]
    H --> I[close tmp]
    I --> J[Root.Rename name.tmp → name]
    J --> K[chmod dst to src perms]
    K --> L[UtimesNano dst atime/mtime]
    L --> M[done]
```

### Guarantees

| Property                     | How                                                                                  |
|------------------------------|--------------------------------------------------------------------------------------|
| Only regular files           | `os.Lstat` + `Mode().IsRegular()`, so symlinks, devices and FIFOs are rejected        |
| Source confinement           | The relative path must stay inside `srcRoot`. It is opened through `os.Root`, which blocks `..` and symlink escapes |
| Destination confinement      | The temp file and the rename both go through an `os.Root` opened on the destination directory |
| Atomic replacement           | Data is written to `<name>.tmp`, `fsync`'d, closed, then **renamed** over `<name>`    |
| Durability                   | `out.Sync()` runs before the rename                                                  |
| Permissions preserved        | The temp file is created with the source permission bits, then `chmod` re-applies them after the rename |
| Timestamps preserved         | Source `atime`/`mtime` are read with `unix.Stat` and restored with `unix.UtimesNano` (nanosecond precision) |
| Efficient I/O                | 1 MiB copy buffer                                                                    |

### Overwrite semantics

An existing `<destination>/<name>` is **silently replaced**, and so is a leftover `<name>.tmp`. There is no
collision handling yet; the code has a `TODO: apply naming scheme with increment`.

## Progress reporting

```go
type ProgressFn func(copied, total int64)
```

- `progressWriter` wraps the destination writer and adds the bytes written to an `atomic.Int64`.
- `startProgressWriter` starts a goroutine with a **150 ms ticker** that calls `progress(copied, total)`. It
  makes one last call when it is stopped, so 100% is always reported. Passing a `nil` callback turns progress off.
- In the TUI, the callback sends `copyProgressMsg{Index, Copied, Total}` on an unbuffered channel. `copyDoneMsg`
  marks a file as complete. `copyFinishedMsg` moves the UI to the `finished` step, then Leafy unmounts the devices
  and quits automatically.
- The view shows each file as `name copied/total (NN%)`.

## Error handling

| Failure point               | Result                                                                        |
|-----------------------------|-------------------------------------------------------------------------------|
| Any step of a single copy   | Returns a wrapped error (`stat src: …`, `copy: …`, `rename dst: …`, etc.)     |
| In the TUI                  | The worker sends `copyErrorMsg` and **stops**. Remaining files are skipped, the error is shown, and Leafy unmounts the devices and quits automatically |
| Already-copied files        | Stay at the destination                                                       |

See [Known Limitations](../known-limitations.md) for edge cases in error cleanup.
