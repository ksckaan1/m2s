# m2s (multipart form to struct converter)

## Installation
```sh
go get -u github.com/ksckaan1/m2s@latest
```

## Usage

```go
type MyRequestBody struct {
  Name          string                `form:"name"`
  Age           int                   `form:"age"`
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