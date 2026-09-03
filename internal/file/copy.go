package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sys/unix"
)

// CopyWithProgress copies a file from src to dst while providing progress updates via the progress callback.
func CopyWithProgress(srcRoot, src, dst string, progress ProgressFn) error {
	// reserve permission
	si, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("stat src: %w", err)
	}
	if !si.Mode().IsRegular() {
		return fmt.Errorf("src is not a regular file: %s (%v)", src, si.Mode())
	}
	total := si.Size()

	// reserve timestamps
	atime, mtime, err := statATimeMTime(src)
	if err != nil {
		return fmt.Errorf("get atime and mtime: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return fmt.Errorf("mkdir dst dir: %w", err)
	}

	srcRel, err := filepath.Rel(srcRoot, src)
	if err != nil {
		return fmt.Errorf("make source relative: %w", err)
	}
	if srcRel == "." || srcRel == ".." || strings.HasPrefix(srcRel, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("source path escapes root: %s", src)
	}

	srcRootHandle, err := os.OpenRoot(srcRoot)
	if err != nil {
		return fmt.Errorf("open source root: %w", err)
	}
	defer srcRootHandle.Close()

	in, err := srcRootHandle.Open(srcRel)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)

	// write to a temporary file
	dstDir := filepath.Dir(dst)
	dstName := filepath.Base(dst)
	if dstName == "." || dstName == "" {
		return fmt.Errorf("invalid destination file name: %q", dst)
	}

	root, err := os.OpenRoot(dstDir)
	if err != nil {
		return fmt.Errorf("open destination root: %w", err)
	}
	defer root.Close()

	// write to a temporary file
	tmp := dstName + ".tmp"
	out, err := root.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, si.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}

	var copied atomic.Int64
	stopProgress := startProgressWriter(&copied, total, progress, 150*time.Millisecond)

	pw := &progressWriter{
		w:      out,
		copied: &copied,
	}
	buf := make([]byte, 1024*1024)
	_, copyErr := io.CopyBuffer(pw, in, buf)

	stopProgress()

	syncErr := out.Sync()

	if copyErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("copy: %w", copyErr)
	}
	if syncErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("sync: %w", syncErr)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("close dst: %w", err)
	}

	if err := root.Rename(tmp, dstName); err != nil {
		return fmt.Errorf("rename dst: %w", err)
	}

	if err := os.Chmod(dst, si.Mode().Perm()); err != nil {
		return fmt.Errorf("chmod dst: %w", err)
	}
	if err := setATimeMTime(dst, atime, mtime); err != nil {
		return fmt.Errorf("set dst timestamps: %w", err)
	}

	return nil
}

// statATimeMTime retrieves access and modification timestamps of the file at the specified path.
func statATimeMTime(path string) (atime, mtime unix.Timespec, err error) {
	var st unix.Stat_t
	if err := unix.Stat(path, &st); err != nil {
		return unix.Timespec{}, unix.Timespec{}, err
	}
	return st.Atim, st.Mtim, nil
}

// setATimeMTime sets the access and modification timestamps of the file at the specified path.
func setATimeMTime(path string, atime, mtime unix.Timespec) error {
	return unix.UtimesNano(path, []unix.Timespec{atime, mtime})
}
