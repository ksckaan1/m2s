# m2s (multipart form to struct converter)

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