package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	dev "github.com/nightails/leafy/internal/device"
	"github.com/nightails/leafy/internal/file"
)

// initDevicesCmd finds and mounts usb devices.
func initDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		var mdevs []device
		devs, err := dev.FindUSBDevices()
		if err != nil {
			return errMsg(err)
		}
		for _, d := range devs {
			if d.Mountpoint == "" {
				d, err = dev.MountDevice(d)
				if err != nil {
					return errMsg(err)
				}
				mdevs = append(mdevs, device{
					d.Name,
					d.Path,
					d.Mountpoint,
				})
			} else {
				mdevs = append(mdevs, device{
					d.Name,
					d.Path,
					d.Mountpoint,
				})
			}
		}
		return devicesMsg(mdevs)
	}
}

// removeDevicesCmd unmounts all usb devices.
func removeDevicesCmd(devs []device) tea.Cmd {
	return func() tea.Msg {
		if len(devs) == 0 {
			return nil
		}

		for _, d := range devs {
			if _, err := dev.UnmountDevice(dev.USBDevice{
				Name:       d.name,
				Path:       d.path,
				Mountpoint: d.mountpoint,
			}); err != nil {
				return errMsg(err)
			}
		}
		return nil
	}
}

// findMediaCmd searches for supported file formats and returns a mediaList.
func findMediaCmd(devices []device) tea.Cmd {
	return func() tea.Msg {
		if len(devices) == 0 {
			return nil
		}

		var paths []string
		for _, d := range devices {
			paths = append(paths, d.mountpoint)
		}

		var media []medium
		files, err := file.GetFiles(paths, formats)
		if err != nil {
			return errMsg(err)
		}
		if len(files) == 0 {
			return nil
		}
		for _, f := range files {
			media = append(media, medium{
				name:   f.Name,
				format: f.Ext,
				root:   f.Root,
				src:    f.Path,
				total:  f.Size,
			})
		}
		return mediaMsg(media)
	}
}

// copyMediaCmd copies given media to the destination path.
func copyMediaCmd(media []medium) tea.Cmd {
	return func() tea.Msg {
		if len(media) == 0 {
			return nil
		}

		ch := make(chan tea.Msg)

		go func() {
			defer close(ch)

			for i := range media {
				if err := file.CopyWithProgress(media[i].root, media[i].src, media[i].dest, func(copied, total int64) {
					ch <- copyProgressMsg{i, copied, total}
				}); err != nil {
					ch <- copyErrorMsg{i, err}
					return
				}

				ch <- copyDoneMsg{i}
			}

			ch <- copyFinishedMsg{}
		}()

		return copyStartedMsg{ch}
	}
}

func listenCopyCmd(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return copyFinishedMsg{}
		}
		return msg
	}
}
