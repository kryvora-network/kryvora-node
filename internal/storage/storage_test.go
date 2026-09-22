package storage

import (
	"os"
	"testing"
)

func TestStorageOpenClose(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kryvora-storage-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := Open(tempDir)
	if err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}

	if err := store.Sync(); err != nil {
		t.Errorf("sync returned error: %v", err)
	}

	if err := store.Close(); err != nil {
		t.Errorf("close returned error: %v", err)
	}

	expected := tempDir + "/test.db"
	if store.Path("test.db") != expected {
		t.Errorf("expected path %s, got %s", expected, store.Path("test.db"))
	}
}
