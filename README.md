# m2s (multipart form to struct converter)

[![tag](https://img.shields.io/github/tag/ksckaan1/m2s.svg)](https://github.com/ksckaan1/m2s/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23.4-%23007d9c)
[![Go report](https://goreportcard.com/badge/github.com/ksckaan1/m2s)](https://goreportcard.com/report/github.com/ksckaan1/m2s)
![m2s](https://img.shields.io/badge/coverage-100%25-green?style=flat)
[![Contributors](https://img.shields.io/github/contributors/ksckaan1/m2s)](https://github.com/ksckaan1/m2s/graphs/contributors)
[![License](https://img.shields.io/github/license/ksckaan1/m2s)](./LICENSE)

## Installation
```sh
go get -u github.com/ksckaan1/m2s@latest
```

## Example Usage

```go
type MyRequestBody struct {
  Name          string                `form:"name"`
  Age           int                   `form:"age"`
  Hobbies       []string              `form:"hobbies"`
  ProfilePhoto  *multipart.FileHeader `form:"profile_photo"`
}
// convert 
var myRequestBody MyRequestBody
// mpf is a *multipart.Form
err := m2s.Convert(mpf, &myRequestBody)
if err != nil {
  return err
}
```

## Supported Types

### Form Files
- `multipart.FileHeader`, `*multipart.FileHeader` for single file
- `[]multipart.FileHeader`, `[]*multipart.FileHeader` for multiple files

### Form Values
- `string`, its pointer and all types derived from its
- All `int` types, its pointers and all types derived from its
- All `uint` types, its pointers and all types derived from its
- All `float` types, its pointers and all types derived from its
- All `bool` types, its pointers and all types derived from its
- All `complex` types, its pointers and all types derived from its
- `struct`, `*struct` (default: json decode)
- Slices of all supported types (default: json decode)
- Maps of all supported types (default: json decode)

> [!NOTE]  
> If field type implements `encoding.TextUnmarshaler`, decodes this field using `UnmarshalText` method.
> 
> Otherwise, decodes this field using `encoding/json.Unmarshal` or parsing primitive value.