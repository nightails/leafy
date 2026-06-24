package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Run("Delete existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "testfile")
		if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Delete(tmpFile); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
			t.Errorf("expected file to be deleted, but it still exists or error: %v", err)
		}
	})

	t.Run("Delete non-existing file (idempotent)", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "nonexistent")

		if err := Delete(tmpFile); err != nil {
			t.Errorf("expected no error for non-existing file, got %v", err)
		}
	})

	t.Run("Error when deleting directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		subDir := filepath.Join(tmpDir, "subdir")
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := Delete(subDir); err == nil {
			t.Error("expected error when deleting a directory, got nil")
		}
	})

	t.Run("Handle error from os.Remove", func(t *testing.T) {
		// To simulate an error from os.Remove, we can try to delete a file
		// in a directory where we don't have write permissions.
		// However, on some systems/environments, this might be tricky to set up reliably in a test.
		// Alternatively, we can use a read-only directory.
		
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "readonly", "file")
		if err := os.Mkdir(filepath.Join(tmpDir, "readonly"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
			t.Fatal(err)
		}
		
		// Make the directory read-only to prevent deletion of files inside it
		if err := os.Chmod(filepath.Join(tmpDir, "readonly"), 0555); err != nil {
			t.Fatal(err)
		}
		defer os.Chmod(filepath.Join(tmpDir, "readonly"), 0755) // cleanup

		err := Delete(tmpFile)
		if err == nil {
			// In the current implementation, this will FAIL (it returns nil despite os.Remove failing)
			// This test case demonstrates the bug.
			t.Error("expected error when os.Remove fails, but got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}
