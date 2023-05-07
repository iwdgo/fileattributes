// Package fileattributes has methods to work on windows file attributes.

//go:build windows
// +build windows

package fileattributes

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	FILE_ATTRIBUTE_READONLY              = 0x00000001
	FILE_ATTRIBUTE_HIDDEN                = 0x00000002
	FILE_ATTRIBUTE_SYSTEM                = 0x00000004
	FILE_ATTRIBUTE_DIRECTORY             = 0x00000010
	FILE_ATTRIBUTE_ARCHIVE               = 0x00000020
	FILE_ATTRIBUTE_DEVICE                = 0x00000040 // System use
	FILE_ATTRIBUTE_NORMAL                = 0x00000080
	FILE_ATTRIBUTE_TEMPORARY             = 0x00000100
	FILE_ATTRIBUTE_SPARSE_FILE           = 0x00000200
	FILE_ATTRIBUTE_REPARSE_POINT         = 0x00000400
	FILE_ATTRIBUTE_OFFLINE               = 0x00001000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED   = 0x00002000
	FILE_ATTRIBUTE_ENCRYPTED             = 0x00004000
	FILE_ATTRIBUTE_INTEGRITY_STREAM      = 0x00008000
	FILE_ATTRIBUTE_VIRTUAL               = 0x00010000 // System use
	FILE_ATTRIBUTE_NO_SCRUB_DATA         = 0x00020000
	FILE_ATTRIBUTE_RECALL_ON_OPEN        = 0x00040000
	FILE_ATTRIBUTE_PINNED                = 0x00080000
	FILE_ATTRIBUTE_UNPINNED              = 0x00100000
	FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS = 0x00400000
)

type FileAttributes uint32

// StatFileAttributes returns the attributes of a path.
// Several calls are attempted until some attributes are available.
// Error returned is always from CreateFile and attributes if any.
func StatFileAttributes(path string) (fa FileAttributes, err error) {
	if fa, err = GetFileAttributesEx(path); err == nil {
		if fa&FILE_ATTRIBUTE_NORMAL == 0 {
			return fa, nil
		}
	}
	if fa, err = FindFirstFile(path); err == nil {
		if fa&FILE_ATTRIBUTE_NORMAL == 0 {
			return fa, nil
		}
	}
	if fa, err = CreateFile(path); err == nil {
		return fa, nil
	}
	return fa, err
}

// GetFileAttributesEx returns attributes of file using GetFileAttributesEx call of Win32
func GetFileAttributesEx(path string) (FileAttributes, error) {
	var fa syscall.Win32FileAttributeData
	namep, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return FileAttributes(0), err
	}
	err = syscall.GetFileAttributesEx(namep, syscall.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&fa)))
	return FileAttributes(fa.FileAttributes), err
}

// FindFirstFile returns attributes of file using FindFirstFile call of Win32
func FindFirstFile(path string) (FileAttributes, error) {
	var fd syscall.Win32finddata
	namep, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return FileAttributes(0), err
	}
	sh, err := syscall.FindFirstFile(namep, &fd)
	if err != nil {
		return FileAttributes(0), err
	}
	defer func() {
		err = syscall.FindClose(sh)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return FileAttributes(fd.FileAttributes), nil
}

// CreateFile acquires a handle of an existing file. The handle is used by GetFileInformationByHandle
func CreateFile(path string) (FileAttributes, error) {
	var d syscall.ByHandleFileInformation
	namep, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return FileAttributes(0), err
	}

	// TODO Does not follow symlink
	h, err := syscall.CreateFile(namep, 0, 0, nil, syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return FileAttributes(0), err
	}
	defer func() {
		err = syscall.CloseHandle(h)
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = syscall.GetFileInformationByHandle(h, &d)
	return FileAttributes(d.FileAttributes), err
}

// PrintAttributes prints attributes using readable names from documentation to os.Stdout when no file is provided.
func PrintAttributes(attrs FileAttributes, f ...*os.File) {
	w := os.Stdout
	if f != nil {
		w = f[0]
	}
	printBit := func(s string, b FileAttributes) {
		if b != 0 {
			_, _ = fmt.Fprintf(w, " %s", s)
		}
	}
	printBit("READONLY", attrs&FILE_ATTRIBUTE_READONLY)
	printBit("HIDDEN", attrs&FILE_ATTRIBUTE_HIDDEN)
	printBit("SYSTEM", attrs&FILE_ATTRIBUTE_SYSTEM)
	printBit("DIRECTORY", attrs&FILE_ATTRIBUTE_DIRECTORY)
	printBit("ARCHIVE", attrs&FILE_ATTRIBUTE_ARCHIVE)
	printBit("DEVICE", attrs&FILE_ATTRIBUTE_DEVICE)
	printBit("NORMAL", attrs&FILE_ATTRIBUTE_NORMAL)
	printBit("TEMPORARY", attrs&FILE_ATTRIBUTE_TEMPORARY)
	printBit("SPARSE_FILE", attrs&FILE_ATTRIBUTE_SPARSE_FILE)
	printBit("REPARSE_POINT", attrs&FILE_ATTRIBUTE_REPARSE_POINT)
	printBit("OFFLINE", attrs&FILE_ATTRIBUTE_OFFLINE)
	printBit("NOT_CONTENT_INDEXED", attrs&FILE_ATTRIBUTE_NOT_CONTENT_INDEXED)
	printBit("ENCRYPTED", attrs&FILE_ATTRIBUTE_ENCRYPTED)
	printBit("INTEGRITY_STREAM", attrs&FILE_ATTRIBUTE_INTEGRITY_STREAM)
	printBit("VIRTUAL", attrs&FILE_ATTRIBUTE_VIRTUAL)
	printBit("NO_SCRUB_DATA", attrs&FILE_ATTRIBUTE_NO_SCRUB_DATA)
	printBit("RECALL_ON_OPEN", attrs&FILE_ATTRIBUTE_RECALL_ON_OPEN)
	printBit("PINNED", attrs&FILE_ATTRIBUTE_PINNED)
	printBit("UNPINNED", attrs&FILE_ATTRIBUTE_UNPINNED)
	printBit("ON_DATA_ACCESS", attrs&FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS)
	fmt.Print("\n")
}
