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
	`C:\pagefile.sys`,
	`C:\Dumpstack.log`,
	`C:\Dumpstack.log.tmp`,
	"CONIN$",    // os.Stdin.Name() = /dev/stdin
	"CONOUT$",   // os.Stdout.Name() = /dev/stdout
	"link.hard", // hard link to go.mod
	"link.dir",  // directory junction (link) to parent directory of module
	"CON",
	"NUL",
	// link using /D is missing as it requires privileges
}

func ExampleStatFileAttributes() {
	donotexist, fail := 0, 0
	for _, s := range files {
		f, err := fileattributes.FindFirstFile(s)
		switch {
		case err == nil:
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
		case os.IsNotExist(err):
			donotexist++
		default:
			fail++
		}
	}
	if donotexist != 0 {
		fmt.Printf("%d files do not exist\n", donotexist)
	}
	if fail != 0 {
		fmt.Print(fmt.Sprintf(pf, "StatFileAttributes", fail))
	}
	// Output: *.go: ARCHIVE
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// C:\pagefile.sys: HIDDEN SYSTEM ARCHIVE
	// C:\Dumpstack.log: HIDDEN SYSTEM ARCHIVE
	// C:\Dumpstack.log.tmp: HIDDEN SYSTEM ARCHIVE
	// CONIN$: ARCHIVE
	// CONOUT$: ARCHIVE
	// link.hard: ARCHIVE
	// link.dir: DIRECTORY REPARSE_POINT
	// CON: ARCHIVE
	// NUL: ARCHIVE
	// 2 files do not exist
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
			perm++
		case os.IsTimeout(err):
			timeout++
		case errors.Is(err, ERROR_SHARING_VIOLATION):
			fmt.Printf("%s: %s\n", s, err) // 2
		default:
			fmt.Printf("%s: %s\n", s, err)
			fail++
		}
	}
	if donotexist != 0 {
		fmt.Printf("%d files do not exist\n", donotexist)
	}
	if perm != 0 {
		fmt.Printf("Access is denied to %d files\n", perm)
	}
	if timeout != 0 {
		fmt.Printf("GetFileAttributesEx timed out for %d files\n", timeout)
	}
	fmt.Print(fmt.Sprintf(pf, "GetFileAttributesEx", fail))
	// Output:
	// *.go: The filename, directory name, or volume label syntax is incorrect.
	// ..\: DIRECTORY
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// C:\: HIDDEN SYSTEM DIRECTORY
	// C:\pagefile.sys: The process cannot access the file because it is being used by another process.
	// C:\Dumpstack.log: HIDDEN SYSTEM ARCHIVE
	// C:\Dumpstack.log.tmp: The process cannot access the file because it is being used by another process.
	// CONIN$: Incorrect function.
	// CONOUT$: Incorrect function.
	// link.hard: ARCHIVE
	// link.dir: DIRECTORY REPARSE_POINT
	// CON: The parameter is incorrect.
	// NUL: The parameter is incorrect.
	// GetFileAttributesEx fails for 5 files
}

// ExampleFindFirstFile is using the test files to demonstrate usage.
func ExampleFindFirstFile() {
	donotexist, fail := 0, 0
	for _, s := range files {
		f, err := fileattributes.FindFirstFile(s)
		switch {
		case err == nil:
			fmt.Printf("%s:", s)
			fileattributes.PrintAttributes(f)
		case os.IsNotExist(err):
			donotexist++
		default:
			fail++
		}
	}
	if donotexist != 0 {
		fmt.Printf("%d files do not exist\n", donotexist)
	}
	if fail != 0 {
		fmt.Print(fmt.Sprintf(pf, "FindFirstFile", fail))
	}
	// Output:
	// *.go: ARCHIVE
	// \\.\\pipe\trkwks: NORMAL
	// go.mod: ARCHIVE
	// C:\pagefile.sys: HIDDEN SYSTEM ARCHIVE
	// C:\Dumpstack.log: HIDDEN SYSTEM ARCHIVE
	// C:\Dumpstack.log.tmp: HIDDEN SYSTEM ARCHIVE
	// CONIN$: ARCHIVE
	// CONOUT$: ARCHIVE
	// link.hard: ARCHIVE
	// link.dir: DIRECTORY REPARSE_POINT
	// CON: ARCHIVE
	// NUL: ARCHIVE
	// 2 files do not exist
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
	// link.hard: ARCHIVE
	// link.dir: DIRECTORY
	// CreateFile fails for 8 files
}
