# Device Management

Package: `internal/device` (`device.go`, `lsblk.go`, `udisksctl.go`). Used by `initDevicesCmd` and
`removeDevicesCmd` in `app/tui/cmd.go`.

## Data model

```go
type USBDevice struct {
    Name       string // the Model of the device (lsblk MODEL of the parent disk)
    Path       string // the partition path, e.g. /dev/sdb1
    Label      string // filesystem label of the partition
    Mountpoint string // current mountpoint, "" if not mounted
}
```

Each `USBDevice` stands for one **mountable filesystem**, usually a partition, not a physical disk.

## Detection: `FindUSBDevices()`

Leafy runs:

```bash
lsblk -J --tree -f -o NAME,PATH,LABEL,FSTYPE,TRAN,TYPE,MODEL,MOUNTPOINTS
```

It unmarshals the JSON into a tree of `blockDevice` and then filters it:

1. Top-level entries are kept only when `type == "disk"` **and** `tran == "usb"`. Internal SATA/NVMe disks, loop
   devices, and optical drives are skipped.
2. If the disk has **children** (partitions), each child with a non-empty `fstype` becomes a `USBDevice`. `Name` comes
   from the parent disk's `model`.
3. If the disk has **no children**, the disk itself is used when it has a filesystem. This covers "superfloppy"
   layouts, for example some cameras and thumb drives formatted without a partition table.
4. The mountpoint is the first entry of `mountpoints` that is not empty and not `"null"`.

If `lsblk` cannot run or returns invalid JSON, the error goes back to the UI as `errMsg`.

## Mounting: `MountDevice(d)`

- If `d.Mountpoint` is already set, Leafy does nothing.
- Otherwise it runs `udisksctl mount -b <path>` with `LC_ALL=C`:
  - On success, it parses `Mounted /dev/sdb1 at /run/media/user/LABEL.` to get the mountpoint.
  - On failure, if the output contains `... already mounted at <mp>.`, that mountpoint is used and no error is
    raised. Otherwise the error includes the full `udisksctl` output.
  - Output that does not match either pattern raises an `unexpected udiskctl output` error.

`initDevicesCmd` mounts the partitions **one at a time**, and the first failure aborts startup with an error.

## Unmounting: `UnmountDevice(d)`

- If `d.Mountpoint` is empty, Leafy does nothing.
- Otherwise it runs `udisksctl unmount -b <path>` with `LC_ALL=C`. Output containing `is not mounted` counts as
  success.
- On exit, `removeDevicesCmd` unmounts **every** device in `state.devices`, which means every detected USB partition,
  whether Leafy mounted it or not. The first failure stops the loop.
- Exit happens either automatically when the transfer completes or fails, or when the user presses `ctrl+c`. On the
  automatic exit, unmounting finishes before the program quits, so an unmount error is shown on the final screen.

## Power-off (not yet used)

`PowerOffDevice(d)` unmounts the device if needed, then runs `udisksctl power-off -b <path>` to spin it down and detach
it safely. The TUI does not call it yet.

## Security

All `udisksctl` calls first check the path against:

```
^/dev/[A-Za-z0-9._/-]+$
```

Paths that do not match are rejected with `invalid block device path`. This prevents argument or command injection
through crafted device names, and it is why the `exec.Command` calls carry `#nosec G204` annotations for `gosec`.

## Requirements and permissions

- `lsblk` and `udisksctl` must be on `PATH`.
- The user must be allowed to mount removable media through polkit. This is the default for local, active desktop
  sessions. Over SSH or in headless sessions, `udisksctl` may ask for authentication or fail.
