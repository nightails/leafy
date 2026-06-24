// Package device provides functions for managing USB devices.
package device

type USBDevice struct {
	Name       string // the Model of the device
	Path       string // the partition path of the device
	Label      string // the label of the partition
	Mountpoint string // the mountpoint of the partition
}

// FindUSBDevices returns a list of USB devices by running lsblk.
func FindUSBDevices() ([]USBDevice, error) {
	blk, err := readLSBLK()
	if err != nil {
		return nil, err
	}

	devs := make([]USBDevice, 0)
	for _, bd := range blk.BlockDevices {
		// filter out non-USB disks
		if bd.Type == "disk" && bd.Tran == "usb" {
			if len(bd.Children) > 0 {
				for _, child := range bd.Children {
					devs = append(devs, USBDevice{
						Name:       bd.Model,
						Path:       child.Path,
						Label:      child.Label,
						Mountpoint: getMountpoint(child),
					})
				}
			} else {
				devs = append(devs, USBDevice{
					Name:       bd.Model,
					Path:       bd.Path,
					Label:      bd.Label,
					Mountpoint: getMountpoint(bd),
				})
			}
		}
	}

	return devs, nil
}

func getMountpoint(bd blockDevice) string {
	for _, mp := range bd.Mountpoints {
		if mp != "" && mp != "null" {
			return mp
		}
	}
	return ""
}

// MountDevice mounts the given USB device using udisksctl and returns the updated device.
func MountDevice(d USBDevice) (USBDevice, error) {
	// already mounted, skip
	if d.Mountpoint != "" {
		return d, nil
	}

	mp, err := mountUdisks(d.Path)
	if err != nil {
		return d, err
	}
	d.Mountpoint = mp
	return d, nil
}

// UnmountDevice unmounts the given USB device using udisksctl and returns the updated device.
func UnmountDevice(d USBDevice) (USBDevice, error) {
	// not mounted, skip
	if d.Mountpoint == "" {
		return d, nil
	}

	if err := unmountUdisks(d.Path); err != nil {
		return d, err
	}
	d.Mountpoint = ""
	return d, nil
}

// PowerOffDevice powers off the given USB device using udisksctl.
func PowerOffDevice(d USBDevice) error {
	if d.Mountpoint != "" {
		if err := unmountUdisks(d.Path); err != nil {
			return err
		}
	}
	return powerOffUdisks(d.Path)
}
