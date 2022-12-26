[![Go Reference](https://pkg.go.dev/badge/github.com/iwdgo/fileattributes.svg)](https://pkg.go.dev/github.com/iwdgo/windowsattributes)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwdgo/fileattributes)](https://goreportcard.com/report/github.com/iwdgo/windowsattributes)
[![codecov](https://codecov.io/gh/iwdgo/fileattributes/branch/master/graph/badge.svg)](https://codecov.io/gh/iwdgo/windowsattributes)

![GitHub](https://github.com/iwdgo/fileattributes/workflows/GitHub/badge.svg)

# File Attributes on Windows

The set of file attributes can change depending on the used Win32 API call.
[File attributes](https://docs.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants) provides some methods
to detail them.

### Go lang

Some useful documentation links:
- [Documentation of package syscall](https://pkg.go.dev/syscall)  
- [Using a module](https://go.dev/ref/mod)  

Some related issues:
- https://go.dev/issue/25923  
- https://go.dev/issue/41755  

## Some related Windows documentation

https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfileattributesexw  
https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-findfirstfileexw  
https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew  
https://learn.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-win32_file_attribute_data    
https://learn.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants  

### Version history

`v0.1.0` Add API for comparison purposes  
