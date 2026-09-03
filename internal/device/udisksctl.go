package device

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var validBlockDevice = regexp.MustCompile(`^/dev/[A-Za-z0-9._/-]+$`)

func validateBlockDevice(device string) error {
	if !validBlockDevice.MatchString(device) {
		return fmt.Errorf("invalid block device path: %q", device)
	}
	return nil
}

// mountUdisks mounts the given device using udisksctl and returns the mountpoint.
func mountUdisks(device string) (string, error) {
	if err := validateBlockDevice(device); err != nil {
		return "", err
	}

	// matches: Mounted /dev/sda1 at /run/media/user/label.
	var reMounted = regexp.MustCompile(`(?i)Mounted\s+(/dev/\S+)\s+at\s+(.+?)\.?\s*$`)
	// matches: Error mounting /dev/sda1: ... already mounted at /run/media/user/label.
	var reAlreadyMounted = regexp.MustCompile(`(?i)already mounted at\s+(.+?)\.?\s*$`)

	// #nosec G204 -- device path is validated before being passed to udisksctl.
	cmd := exec.Command("udisksctl", "mount", "-b", device)
	cmd.Env = append(os.Environ(), "LC_ALL=C")

	out, err := cmd.CombinedOutput()
	if err != nil {
		// check if already mounted
		m := reAlreadyMounted.FindStringSubmatch(string(out))
		if len(m) >= 2 {
			return strings.TrimSpace(m[1]), nil
		}
		return "", fmt.Errorf("%v: %s", err, string(out))
	}

	m := reMounted.FindStringSubmatch(string(out))
	if len(m) < 3 {
		return "", fmt.Errorf("unexpected udiskctl output: %s", string(out))
	}
	return strings.TrimSpace(m[2]), nil
}

// unmountUdisks unmounts the given device using udisksctl.
func unmountUdisks(device string) error {
	if err := validateBlockDevice(device); err != nil {
		return err
	}

	var reNotMounted = regexp.MustCompile(`(?i)is not mounted`)

	// #nosec G204 -- device path is validated before being passed to udisksctl.
	cmd := exec.Command("udisksctl", "unmount", "-b", device)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if reNotMounted.MatchString(string(out)) {
			return nil
		}
		return fmt.Errorf("%v: %s", err, string(out))
	}
	return nil
}

// powerOffUdisks safely powers off the given device using udisksctl.
func powerOffUdisks(device string) error {
	if err := validateBlockDevice(device); err != nil {
		return err
	}

	// #nosec G204 -- device path is validated before being passed to udisksctl.
	cmd := exec.Command("udisksctl", "power-off", "-b", device)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	_, err := cmd.CombinedOutput()
	return err
}
