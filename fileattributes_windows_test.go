//go:build windows
// +build windows

package fileattributes

import (
	"os"
	"syscall"
	"testing"
)

const (
	archivePath                   = "go.mod"
	pipePath                      = `\\.\\pipe\trkwks`
	ERROR_BUSY      syscall.Errno = 170
	ERROR_PIPE_BUSY syscall.Errno = 231
)

func TestFileArchive(t *testing.T) {
	fa1, err := GetFileAttributesEx(archivePath)
	if err != nil {
		t.Fatalf("%s", err)
	}

	fa2, err := FindFirstFile(archivePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if fa1 != fa2 {
		t.Errorf("got (FindFirstFile) %b, want (GetFileAttributesEx) %b", fa2, fa1)
	}

	fa3, err := CreateFile(archivePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if fa1 != fa3 {
		t.Errorf("got (CreateFile) %b, want (GetFileAttributesEx) %b", fa3, fa1)
	}
	if fa1 != FILE_ATTRIBUTE_ARCHIVE {
		t.Fatalf("got %b (%v), want %b", fa1, fa1, FILE_ATTRIBUTE_ARCHIVE)
	}
}

func TestPipe(t *testing.T) {
	fa1, err := GetFileAttributesEx(pipePath)
	if err != nil {
		t.Fatalf("%s", err)
	}

	fa2, err := FindFirstFile(pipePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if fa1 != fa2 {
		t.Errorf("got (GetFileAttributesEx) %b, want (FindFirstFile) %b", fa1, fa2)
	}

	fa3, err := CreateFile(pipePath)
	switch err {
	case ERROR_PIPE_BUSY:
	case ERROR_BUSY:
		t.Skipf("%s. Ignoring test using CreateFile on a pipe", err)
	default:
		if err != nil {
			t.Fatalf("%s", err)
		}
	}
	if fa1 != fa3 {
		PrintAttributes(fa1)
		t.Errorf("got (GetFileAttributesEx) %b, want (CreateFile) %b", fa1, fa3)
	}
	// No attributes for a pipe
	if fa1 != FILE_ATTRIBUTE_NORMAL {
		PrintAttributes(fa1)
		t.Fatalf("got %b, want %b", fa1, FILE_ATTRIBUTE_ARCHIVE)
	}
}

func TestStatFileAttributes(t *testing.T) {
	fa, err := StatFileAttributes(pipePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
	// On Windows, Win32 API do not return attributes for a pipe
	if fa&FILE_ATTRIBUTE_NORMAL != 0 {
		return
	}
	PrintAttributes(fa)
	t.Fatalf("FILE_ATTRIBUTE_NORMAL is not set: %b", fa)
}

func TestDoesNotExist(t *testing.T) {
	const doesnotexit = "doesnotexit.txt"
	_, err := GetFileAttributesEx(doesnotexit)
	if err == nil {
		t.Fatalf("%s file is not expected to exist", doesnotexit)
	}
	if !os.IsNotExist(err) {
		t.Fatalf("%s", err)
	}
	_, err = FindFirstFile(doesnotexit)
	if !os.IsNotExist(err) {
		t.Fatalf("%s", err)
	}
	_, err = FindFirstFile("*.go2")
	if !os.IsNotExist(err) {
		t.Fatalf("%s", err)
	}
	_, err = CreateFile(doesnotexit)
	if !os.IsNotExist(err) {
		t.Fatalf("%s", err)
	}
}

// TODO Add benchmark
