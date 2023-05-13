//go:build windows
// +build windows

package fileattributes_test

import (
	"errors"
	"fmt"
	fileattributes "github.com/iwdgo/filesattributes"
	"os"
	"syscall"
)

const ERROR_SHARING_VIOLATION syscall.Errno = 32

var files = []string{
	"*.go", // Not exact but files to find
	`..\`,
	`\\.\\pipe\trkwks`, // Pipe presence and use is recommended
	"go.mod",
	`C:\`,
	// `C:\pagefile.sys`, // Fails locally
	// `C:\Dumpstack.log`, // Fails locally
	// `C:\Dumpstack.log.tmp`, // Fails locally
	"CONIN$",    // os.Stdin.Name() = /dev/stdin
	"CONOUT$",   // os.Stdout.Name() = /dev/stdout
	"link.hard", // TODO Does not exist
	"link.dir",  // TODO Does not exist
	"CON",
	"NUL",
	// link using /D is missing
}

func ExampleStatFileAttributes() {
	fail := 0
	for _, s := range files {
		if f, err := fileattributes.FindFirstFile(s); err == nil {
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
			continue
		}
		fail++
	}
	fmt.Print(fmt.Sprintf(pf, "StatFileAttributes", fail))
	// Output: *.go: ARCHIVE
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// CONIN$: ARCHIVE
	// CONOUT$: ARCHIVE
	// CON: ARCHIVE
	// NUL: ARCHIVE
	// StatFileAttributes fails for 4 files
}

// reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7",
//	"COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

const pf = "%s fails for %d files\n"

// ExampleGetFileAttributesEx is using the test files to demonstrate usage.
func ExampleGetFileAttributesEx() {
	donotexist := 0
	perm := 0
	timeout := 0
	fail := 0
	for _, s := range files {
		f, err := fileattributes.GetFileAttributesEx(s)
		switch {
		case err == nil:
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
		case os.IsNotExist(err):
			donotexist++
		case os.IsPermission(err):
		case os.IsTimeout(err):
			perm++
		case errors.Is(err, ERROR_SHARING_VIOLATION):
			fmt.Printf("%s: %s\n", s, err) // 2
		default:
			fmt.Printf("%s: %s\n", s, err)
			fail++
		}
	}
	fmt.Printf("%d files do not exist\n", donotexist)
	fmt.Printf("Access is denied to %d files\n", perm)
	fmt.Printf("GetFileAttributesEx timed out for %d files\n", timeout)
	fmt.Print(fmt.Sprintf(pf, "GetFileAttributesEx", fail))
	// Output:
	// *.go: The filename, directory name, or volume label syntax is incorrect.
	// ..\: DIRECTORY
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// C:\: HIDDEN SYSTEM DIRECTORY
	// CONIN$: Incorrect function.
	// CONOUT$: Incorrect function.
	// CON: The parameter is incorrect.
	// NUL: The parameter is incorrect.
	// 2 files do not exist
	// Access is denied to 0 files
	// GetFileAttributesEx timed out for 0 files
	// GetFileAttributesEx fails for 5 files
}

// ExampleFindFirstFile is using the test files to demonstrate usage.
func ExampleFindFirstFile() {
	fail := 0
	for _, s := range files {
		if f, err := fileattributes.FindFirstFile(s); err == nil {
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
			continue
		}
		fail++
	}
	fmt.Print(fmt.Sprintf(pf, "FindFirstFile", fail))
	// Output:
	// *.go: ARCHIVE
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// CONIN$: ARCHIVE
	// CONOUT$: ARCHIVE
	// CON: ARCHIVE
	// NUL: ARCHIVE
	// FindFirstFile fails for 4 files
}

// ExampleFindFirstFile is using the test files to demonstrate usage.
func ExampleCreateFile() {
	fail := 0
	for _, s := range files {
		if f, err := fileattributes.CreateFile(s); err == nil {
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
			continue
		}
		fail++
	}
	fmt.Print(fmt.Sprintf(pf, "CreateFile", fail))
	// Output: ..\: DIRECTORY
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// C:\: HIDDEN SYSTEM DIRECTORY
	// CreateFile fails for 7 files
}
