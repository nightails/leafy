package device

import (
	"encoding/json"
	"os/exec"
)

type lsblk struct {
	BlockDevices []blockDevice `json:"blockdevices"`
}
type blockDevice struct {
	Name        string        `json:"name"`
	Path        string        `json:"path"`
	Label       string        `json:"label"`
	Tran        string        `json:"tran"`
	Type        string        `json:"type"`
	Model       string        `json:"model"`
	Mountpoints []string      `json:"mountpoints"`
	Children    []blockDevice `json:"children"`
}

var execCommand = exec.Command

func readLSBLK() (lsblk, error) {
	args := []string{
		"-J",
		"--tree",
		"-f",
		"-o", "NAME,PATH,LABEL,TRAN,TYPE,MODEL,MOUNTPOINTS",
	}
	cmd := execCommand("lsblk", args...)
	raw, err := cmd.Output()
	if err != nil {
		return lsblk{}, err
	}

	var blk lsblk
	if err := json.Unmarshal(raw, &blk); err != nil {
		return lsblk{}, err
	}
	return blk, nil
}
