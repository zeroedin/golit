package fileutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWriteNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	if err := WriteFileAtomic(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q, want %q", got, "hello")
	}

	if runtime.GOOS != "windows" {
		info, _ := os.Stat(path)
		if info.Mode().Perm() != 0644 {
			t.Errorf("permissions = %o, want 0644", info.Mode().Perm())
		}
	}
}

func TestOverwritePreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not support fine-grained file permissions")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "existing.txt")

	if err := os.WriteFile(path, []byte("old"), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := WriteFileAtomic(path, []byte("new"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "new" {
		t.Errorf("content = %q, want %q", got, "new")
	}

	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0755 {
		t.Errorf("permissions = %o, want 0755 (should preserve original)", info.Mode().Perm())
	}
}

func TestNonExistentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-dir", "file.txt")

	err := WriteFileAtomic(path, []byte("data"), 0644)
	if err == nil {
		t.Fatal("expected error for non-existent parent directory")
	}

	// Verify no temp file was leaked in the parent (which doesn't exist)
	parent := filepath.Dir(path)
	if _, statErr := os.Stat(parent); !os.IsNotExist(statErr) {
		t.Errorf("parent directory should not exist, but got: %v", statErr)
	}
}

func TestContentCorrectness(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.bin")

	// Write a non-trivial payload
	data := make([]byte, 1<<16) // 64 KiB
	for i := range data {
		data[i] = byte(i % 251)
	}

	if err := WriteFileAtomic(path, data, 0644); err != nil {
		t.Fatalf("WriteFileAtomic: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(got) != len(data) {
		t.Fatalf("len = %d, want %d", len(got), len(data))
	}
	for i := range data {
		if got[i] != data[i] {
			t.Fatalf("byte %d: got %d, want %d", i, got[i], data[i])
			break
		}
	}
}
