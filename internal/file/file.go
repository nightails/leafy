// Package file provides a File Path utility.
package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type File struct {
	Name string
	Ext  string
	Root string
	Path string
	Size int64
}

// GetFiles scans the provided File paths and returns a list of File paths.
func GetFiles(srcPaths, formats []string) ([]File, error) {
	var supportedFiles []File

	for _, p := range srcPaths {
		if p == "" {
			continue
		}

		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			supportedFiles = addFile(
				supportedFiles,
				File{info.Name(), filepath.Ext(p), filepath.Dir(p), p, info.Size()},
				formats,
			)
			continue
		}

		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			inf, err := d.Info()
			if err != nil {
				return err
			}
			supportedFiles = addFile(
				supportedFiles,
				File{d.Name(), filepath.Ext(path), p, path, inf.Size()},
				formats,
			)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return supportedFiles, nil
}

// addFile adds a File Path to the given list if it's a supported format and not already present.
func addFile(list []File, f File, formats []string) []File {
	if !isSupported(f, formats) {
		return list
	}
	if slices.Contains(list, f) {
		return list
	}
	list = append(list, f)
	return list
}

// isSupported returns true if the given File Path has a supported format.
func isSupported(f File, formats []string) bool {
	return slices.Contains(formats, strings.ToLower(f.Ext))
}
