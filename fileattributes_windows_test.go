//go:build windows
// +build windows

package fileattributes

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

const (
	archivePath                   = "go.mod"
	pipePath                      = `\\.\\pipe\trkwks`
	ERROR_BUSY      syscall.Errno = 170
	ERROR_PIPE_BUSY syscall.Errno = 231
)

func TestMain(m *testing.M) {
	d := "target"
	err := os.Mkdir(d, os.FileMode(600))
	if err != nil {
		panic(err)
	}
	dd := "link.dir"
	cmd := exec.Command("cmd", "/C", "mklink", "/J", dd, "target")
	err = cmd.Run()
	// err = os.Symlink(d, "link.dir") // Requires privilege
	if err != nil {
		out, _ := cmd.CombinedOutput()
		log.Printf("%s", out)
		log.Print(err)
	}
	dl := "link.hard"
	cmd = exec.Command("cmd", "/C", "mklink", "/H", dl, "go.mod")
	err = cmd.Run()
	// err = os.Link("go.mod", dl) // Requires privilege
	if err != nil {
		out, _ := cmd.CombinedOutput()
		log.Printf("%s", out)
		log.Print(err)
	}
	e := m.Run()
	_ = os.Remove(d)
	_ = os.Remove(dd)
	_ = os.Remove(dl)
	os.Exit(e)
}

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
	pipeError(t, err)
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
	pipeError(t, err)
	// On Windows, Win32 API does not return attributes for a pipe
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

func TestPrintBit(t *testing.T) {
	f, err := CreateFile(pipePath)
	pipeError(t, err)
	flog, err := os.CreateTemp("", "")
	if err != nil {
		t.Error(err)
	}
	PrintAttributes(f, flog)
	filename := flog.Name()
	_ = flog.Close() // cannot flush a temporary file on windows
	flog, err = os.Open(filename)
	if err != nil {
		t.Error(f)
	}
	buf, err := io.ReadAll(flog)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) == 0 {
		t.Error("logging file is empty")
		return
	}
	if w := []byte("NORMAL"); bytes.Equal(buf, w) {
		t.Errorf("got %s, want %s", buf, w)
	}
}

func pipeError(t *testing.T, err error) {
	switch err {
	case nil:
		return
	case ERROR_PIPE_BUSY:
	case ERROR_BUSY:
		t.Skipf("%s. Ignoring test using CreateFile on a pipe", err)
	default:
		t.Fatalf("%s", err)
	}
}

// TODO Add benchmark
