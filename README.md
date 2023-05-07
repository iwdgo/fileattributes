[![Go Reference](https://pkg.go.dev/badge/github.com/iwdgo/fileattributes.svg)](https://pkg.go.dev/github.com/iwdgo/fileattributes)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwdgo/fileattributes)](https://goreportcard.com/report/github.com/iwdgo/fileattributes)
[![codecov](https://codecov.io/gh/iwdgo/fileattributes/branch/master/graph/badge.svg)](https://codecov.io/gh/iwdgo/fileattributes)

[![Go](https://github.com/iwdgo/fileattributes/actions/workflows/go.yml/badge.svg)](https://github.com/iwdgo/fileattributes/actions/workflows/go.yml)

# File Attributes on Windows

The set of file attributes can change depending on the used Win32 API call.
[File attributes](https://docs.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants) provides some methods
to detail them.

### Go language

Documentation links:
- [Documentation of package syscall](https://pkg.go.dev/syscall)  
- [Using a module](https://go.dev/ref/mod)  

Related issues:
- https://go.dev/issue/25923  
- https://go.dev/issue/41755  

## Windows documentation

- [GetFileAttributesEx](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfileattributesexw)  
- [FileFirstFileExW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-findfirstfileexw)  
- [CreateFileW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew)  
- [File attributes](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-win32_file_attribute_data)    
- [File attributes constants](https://learn.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants)  
