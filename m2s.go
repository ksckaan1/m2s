package m2s

import (
	"cmp"
	"mime/multipart"
	"reflect"
	"strconv"
)

type fieldType uint

const (
	value fieldType = iota
	values
	file
	files
)

func Convert(mpf *multipart.Form, v any) error {
	rv := reflect.ValueOf(v)

	err := validate(rv)
	if err != nil {
		return err
	}

	rv = rv.Elem()

	for i := range rv.NumField() {
		fieldType := rv.Type().Field(i)
		fieldValue := rv.Field(i)
		tag := fieldType.Tag.Get("form")
		if !fieldType.IsExported() || tag == "-" {
			continue // Skip if struct field is unexported or ignored (-)
		}

		fieldName := cmp.Or(tag, fieldType.Name)

		ft := determineFieldType(fieldType.Type)

		if ft == file || ft == files {
			formFiles, ok := mpf.File[fieldName]
			if !ok || len(formFiles) == 0 {
				continue
			}

			// if single file
			if ft == file {
				setFile(fieldType.Type, fieldValue, formFiles[0])
				continue
			}

			// if multiple files
			setFiles(formFiles, fieldType.Type, fieldValue)
			continue
		}

		formValues, ok := mpf.Value[fieldName]
		if !ok || len(formValues) == 0 {
			continue
		}

		if ft == value { // if single value
			err = convertValue(fieldType.Type, fieldValue, formValues[0])
			if err != nil {
				return err
			}
			continue
		}

		// if multiple values
		err = convertValues(formValues, fieldType.Type, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func setFile(fieldType reflect.Type, fieldValue reflect.Value, formFile *multipart.FileHeader) {
	if formFile == nil {
		return
	}
	if fieldType.Kind() != reflect.Pointer {
		fieldValue.Set(reflect.ValueOf(formFile).Elem())
		return
	}
	fieldValue.Set(reflect.ValueOf(formFile))
}

func setFiles(formFiles []*multipart.FileHeader, fieldType reflect.Type, fieldValue reflect.Value) {
	list := reflect.MakeSlice(fieldType, 0, len(formFiles))
	for i := range formFiles {
		v := reflect.New(reflect.TypeFor[multipart.FileHeader]())
		setFile(v.Type().Elem(), v.Elem(), formFiles[i])
		if fieldType.Elem().Kind() == reflect.Pointer {
			list = reflect.Append(list, v)
		} else {
			list = reflect.Append(list, v.Elem())
		}
	}
	fieldValue.Set(list)
}

func convertValue(fieldType reflect.Type, fieldValue reflect.Value, formValue string) error {
	switch fieldType.Kind() {
	case reflect.Pointer:
		v := reflect.New(fieldType.Elem())
		err := convertValue(fieldType.Elem(), v.Elem(), formValue)
		if err != nil {
			return err
		}
		fieldValue.Set(v)
	case reflect.String:
		fieldValue.SetString(formValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(formValue, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(formValue, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(formValue, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(formValue)
		if err != nil {
			return err
		}
		fieldValue.SetBool(v)
	case reflect.Complex64, reflect.Complex128:
		v, err := strconv.ParseComplex(formValue, 64)
		if err != nil {
			return err
		}
		fieldValue.SetComplex(v)
	default:
		return ErrInvalidFieldType
	}
	return nil
}

func convertValues(formValues []string, fieldType reflect.Type, fieldValue reflect.Value) error {
	list := reflect.MakeSlice(fieldType, 0, len(formValues))
	for i := range formValues {
		v := reflect.New(fieldType.Elem())
		if fieldType.Elem().Kind() == reflect.Pointer {
			v = reflect.New(fieldType.Elem().Elem())
		}
		err := convertValue(v.Type().Elem(), v.Elem(), formValues[i])
		if err != nil {
			return err
		}
		if fieldType.Elem().Kind() == reflect.Pointer {
			list = reflect.Append(list, v)
		} else {
			list = reflect.Append(list, v.Elem())
		}
	}
	fieldValue.Set(list)
	return nil
}

func validate(rv reflect.Value) error {
	if rv.Kind() != reflect.Ptr {
		return ErrValueMustBePointer
	}
	if rv.IsNil() {
		return ErrValueCannotBeNil
	}
	if rv.Elem().Kind() != reflect.Struct {
		return ErrValueMustBeStruct
	}
	return nil
}

func determineFieldType(rt reflect.Type) fieldType {
	if rt.Kind() == reflect.Pointer && rt.Elem() == reflect.TypeFor[multipart.FileHeader]() ||
		rt == reflect.TypeFor[multipart.FileHeader]() {
		return file
	} else if rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Pointer && rt.Elem().Elem() == reflect.TypeFor[multipart.FileHeader]() || // []*multipart.File
		rt.Kind() == reflect.Slice && rt.Elem() == reflect.TypeFor[multipart.FileHeader]() {
		return files
	}
	if rt.Kind() == reflect.Slice {
		return values
	}
	return value
}
