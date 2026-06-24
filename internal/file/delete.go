package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Delete(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("not allowed to delete a directory: %s", path)
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	// Sync parent directory to ensure the deletion is persisted to disk.
	parent := filepath.Dir(path)
	f, err := os.Open(parent)
	if err != nil {
		// If we can't open the parent directory, we just return nil because the file is already removed.
		// This could happen if the parent directory was also removed or due to permissions.
		return nil
	}
	defer f.Close()

	if err := f.Sync(); err != nil {
		// Similar to Open, if Sync fails we don't necessarily want to return an error to the user
		// as the primary goal (deleting the file) has been achieved.
		return nil
	}

	return nil
}
