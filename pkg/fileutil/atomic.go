// Package fileutil provides file I/O utilities.
package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFileAtomic writes data to destPath atomically using a
// temp-file + fsync + rename pattern. If destPath already exists,
// its permissions are preserved; otherwise perm is used.
//
// The temp file is created in the same directory as destPath so that
// the final os.Rename is guaranteed to be atomic on POSIX filesystems.
func WriteFileAtomic(destPath string, data []byte, perm os.FileMode) error {
	// Preserve existing file permissions when overwriting.
	if info, err := os.Stat(destPath); err == nil {
		perm = info.Mode().Perm()
	}

	dir := filepath.Dir(destPath)
	tmp, err := os.CreateTemp(dir, ".golit-tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	success := false
	defer func() {
		if !success {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmp.Chmod(perm); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("syncing temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		d.Close()
	}

	success = true
	return nil
}
