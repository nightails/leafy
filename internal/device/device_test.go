package device

import (
	"os/exec"
	"testing"
)

func TestFindUSBDevices(t *testing.T) {
	// Mock lsblk output
	lsblkJSON := `{
		"blockdevices": [
			{
				"name": "sda",
				"path": "/dev/sda",
				"tran": "usb",
				"type": "disk",
				"model": "USB_Flash_Drive",
				"mountpoints": [null],
				"children": [
					{
						"name": "sda1",
						"path": "/dev/sda1",
						"label": "PART1",
						"mountpoints": ["/mnt/part1"]
					},
					{
						"name": "sda2",
						"path": "/dev/sda2",
						"label": "PART2",
						"mountpoints": [null]
					}
				]
			},
			{
				"name": "sdb",
				"path": "/dev/sdb",
				"tran": "usb",
				"type": "disk",
				"model": "Single_Partition_Disk",
				"mountpoints": ["/mnt/sdb"],
				"children": null
			}
		]
	}`

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		return exec.Command("echo", lsblkJSON)
	}

	devices, err := FindUSBDevices()
	if err != nil {
		t.Fatalf("FindUSBDevices failed: %v", err)
	}

	// Desired implementation is expected to:
	// 1. Find sda1 and sda2 from sda
	// 2. Find sdb even if it has no children

	expectedCount := 3
	if len(devices) != expectedCount {
		t.Fatalf("Expected %d devices, got %d. Devices: %+v", expectedCount, len(devices), devices)
	}

	// Verify sda1
	if devices[0].Path != "/dev/sda1" || devices[0].Label != "PART1" || devices[0].Mountpoint != "/mnt/part1" {
		t.Errorf("Unexpected device 0: %+v", devices[0])
	}

	// Verify sda2
	if devices[1].Path != "/dev/sda2" || devices[1].Label != "PART2" || devices[1].Mountpoint != "" {
		t.Errorf("Unexpected device 1: %+v", devices[1])
	}

	// Verify sdb
	if devices[2].Path != "/dev/sdb" || devices[2].Mountpoint != "/mnt/sdb" {
		t.Errorf("Unexpected device 2: %+v", devices[2])
	}
}

func TestMountDevice_AlreadyMounted(t *testing.T) {
	device := USBDevice{
		Name:  "Test Disk",
		Path:  "/dev/sdc1",
		Label: "TEST",
	}

	udisksOutput := "Error mounting /dev/sdc1: GDBus.Error:org.freedesktop.UDisks2.Error.AlreadyMounted: Device /dev/sdc1 is already mounted at /media/user/TEST."

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "echo '"+udisksOutput+"'; exit 1")
	}

	updated, err := MountDevice(device)
	if err != nil {
		t.Fatalf("MountDevice failed: %v", err)
	}

	expectedMountpoint := "/media/user/TEST"
	if updated.Mountpoint != expectedMountpoint {
		t.Errorf("Expected mountpoint %s, got %s", expectedMountpoint, updated.Mountpoint)
	}
}
