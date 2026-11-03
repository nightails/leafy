# Media Discovery

Code: `internal/file/file.go` (`GetFiles`) and `app/tui/format.go` (supported formats), called from `findMediaCmd`
in `app/tui/cmd.go`.

## Supported formats

| Type  | Extensions                                                  |
|-------|-------------------------------------------------------------|
| Audio | `.mp3`, `.wav`, `.flac`                                     |
| Video | `.mp4`, `.mkv`, `.avi`, `.mov`, `.wmv`, `.flv`, `.webm`     |

Extensions are compared **case-insensitively**. The extension is lower-cased before lookup, so `GX010001.MP4` and
`REC.WAV` match. To add a format, append it, lower-cased and with the leading dot, to `audio` or `video` in
`app/tui/format.go`.

## Scanning algorithm: `GetFiles(srcPaths, formats)`

`srcPaths` holds the mountpoints of all detected USB partitions. For each path:

1. Empty paths are skipped.
2. `os.Stat` the path. An error aborts the whole scan.
3. **If it is a file**, it is checked directly. Its `Root` is the file's parent directory.
4. **If it is a directory**, it is walked recursively with `filepath.WalkDir`:
   - Directories are traversed, and **every** directory is included. Hidden directories and system folders such as
     `DCIM/`, `MISC/`, or `.Trashes/` are not filtered out.
   - Every non-directory entry is checked against the formats. Its `Root` is the scanned mountpoint.
   - Any walk error, such as a permission problem or I/O error, aborts the scan.
5. A file is added only if its extension is supported **and** an identical `File` value is not already in the list.

The result is a slice of:

```go
type File struct {
    Name string // base name, e.g. "GX010001.MP4"
    Ext  string // extension as found, e.g. ".MP4"
    Root string // scan root (mountpoint) used to confine later copying
    Path string // absolute path to the file
    Size int64  // size in bytes
}
```

## Ordering

Devices appear in `lsblk` order. Within a device, files come in `WalkDir` order, which is lexical by path. The UI list
uses the same order.

## Behaviour in the UI

- If no files are found, `findMediaCmd` returns `nil`, and the list stays empty with no extra message.
- Each list row shows `name (size)`, using decimal units (B, KB, MB, GB, TB, with two decimals).

## Caveats

- Symlinks are listed if their name ends in a supported extension. Copying them then fails, because copying only
  accepts regular files.
- macOS "AppleDouble" sidecar files (`._GX010001.MP4`) match the extension check and show up in the list.
- Files with the **same name** in different folders or on different devices are all listed. Because the destination
  layout is flat, they overwrite each other when copied. See [File Transfer](file-transfer.md).
