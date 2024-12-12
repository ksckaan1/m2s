package m2s

import (
	"encoding"
	"errors"
	"mime/multipart"
	"reflect"
	"testing"
	"time"
)

// required for testing, not for normal usage

var _ encoding.TextUnmarshaler = (*failedCustomText)(nil)

type failedCustomText struct{}

func (f *failedCustomText) UnmarshalText(_ []byte) error {
	return errors.New("example error")
}

func TestConvert(t *testing.T) {
	type CustomType struct {
		Name string
		Age  int
	}

	testTime := time.Now().Round(1 * time.Second)

	tests := []struct {
		name             string
		fillMulipartForm func(*multipart.Form) error
		v                any
		wantValue        any
		checkError       func(*testing.T, error)
	}{
		{
			name: "primitive types",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["key1"] = []string{"value1"}
				mpf.Value["key2"] = []string{"42"}
				mpf.Value["key3"] = []string{"42"}
				mpf.Value["key4"] = []string{"42.5"}
				mpf.Value["key5"] = []string{"10+11i"}
				mpf.Value["key6"] = []string{"true"}
				return nil
			},
			v: &struct {
				Key1 string    `form:"key1"`
				Key2 int32     `form:"key2"`
				Key3 uint16    `form:"key3"`
				Key4 float32   `form:"key4"`
				Key5 complex64 `form:"key5"`
				Key6 bool      `form:"key6"`
			}{},
			wantValue: &struct {
				Key1 string    `form:"key1"`
				Key2 int32     `form:"key2"`
				Key3 uint16    `form:"key3"`
				Key4 float32   `form:"key4"`
				Key5 complex64 `form:"key5"`
				Key6 bool      `form:"key6"`
			}{
				Key1: "value1",
				Key2: 42,
				Key3: 42,
				Key4: 42.5,
				Key5: 10 + 11i,
				Key6: true,
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "pointer types",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["key1"] = []string{"value1"}
				mpf.Value["key2"] = []string{"42"}
				mpf.Value["key3"] = []string{"42"}
				mpf.Value["key4"] = []string{"42.5"}
				mpf.Value["key5"] = []string{"10+11i"}
				mpf.Value["key6"] = []string{"true"}
				return nil
			},
			v: &struct {
				Key1 *string    `form:"key1"`
				Key2 *int32     `form:"key2"`
				Key3 *uint16    `form:"key3"`
				Key4 *float32   `form:"key4"`
				Key5 *complex64 `form:"key5"`
				Key6 *bool      `form:"key6"`
			}{},
			wantValue: &struct {
				Key1 *string    `form:"key1"`
				Key2 *int32     `form:"key2"`
				Key3 *uint16    `form:"key3"`
				Key4 *float32   `form:"key4"`
				Key5 *complex64 `form:"key5"`
				Key6 *bool      `form:"key6"`
			}{
				Key1: ptr("value1"),
				Key2: ptr[int32](42),
				Key3: ptr[uint16](42),
				Key4: ptr[float32](42.5),
				Key5: ptr[complex64](10 + 11i),
				Key6: ptr(true),
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "file",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.File["file1"] = []*multipart.FileHeader{
					{
						Filename: "file1.txt",
					},
				}
				mpf.File["file2"] = []*multipart.FileHeader{
					{
						Filename: "file2.txt",
					},
				}
				return nil
			},
			v: &struct {
				File1 *multipart.FileHeader `form:"file1"`
				File2 multipart.FileHeader  `form:"file2"`
			}{},
			wantValue: &struct {
				File1 *multipart.FileHeader `form:"file1"`
				File2 multipart.FileHeader  `form:"file2"`
			}{
				File1: &multipart.FileHeader{Filename: "file1.txt"},
				File2: multipart.FileHeader{Filename: "file2.txt"},
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "file list",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.File["files1"] = []*multipart.FileHeader{
					{
						Filename: "file1.txt",
					},
					{
						Filename: "file2.txt",
					},
				}
				mpf.File["files2"] = []*multipart.FileHeader{
					{
						Filename: "file3.txt",
					},
					{
						Filename: "file4.txt",
					},
				}
				return nil
			},
			v: &struct {
				Files1 []*multipart.FileHeader `form:"files1"`
				Files2 []multipart.FileHeader  `form:"files2"`
			}{},
			wantValue: &struct {
				Files1 []*multipart.FileHeader `form:"files1"`
				Files2 []multipart.FileHeader  `form:"files2"`
			}{
				Files1: []*multipart.FileHeader{
					{Filename: "file1.txt"},
					{Filename: "file2.txt"},
				},
				Files2: []multipart.FileHeader{
					{Filename: "file3.txt"},
					{Filename: "file4.txt"},
				},
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "time type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["time1"] = []string{testTime.Format(time.RFC3339)}
				mpf.Value["time2"] = []string{testTime.Format(time.RFC3339)}
				return nil
			},
			v: &struct {
				Time1 *time.Time `form:"time1"`
				Time2 time.Time  `form:"time2"`
			}{},
			wantValue: &struct {
				Time1 *time.Time `form:"time1"`
				Time2 time.Time  `form:"time2"`
			}{
				Time1: &testTime,
				Time2: testTime,
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "struct type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["custom1"] = []string{`{"name":"John","age":42}`}
				mpf.Value["custom2"] = []string{`{"name":"Jane","age":21}`}
				return nil
			},
			v: &struct {
				Custom1 CustomType  `form:"custom1"`
				Custom2 *CustomType `form:"custom2"`
			}{},
			wantValue: &struct {
				Custom1 CustomType  `form:"custom1"`
				Custom2 *CustomType `form:"custom2"`
			}{
				Custom1: CustomType{
					Name: "John",
					Age:  42,
				},
				Custom2: &CustomType{
					Name: "Jane",
					Age:  21,
				},
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "slice type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["custom1"] = []string{`[{"name":"John","age":42}]`}
				mpf.Value["custom2"] = []string{`[{"name":"Jane","age":21}]`}
				return nil
			},
			v: &struct {
				Custom1 []CustomType  `form:"custom1"`
				Custom2 []*CustomType `form:"custom2"`
			}{},
			wantValue: &struct {
				Custom1 []CustomType  `form:"custom1"`
				Custom2 []*CustomType `form:"custom2"`
			}{
				Custom1: []CustomType{
					{
						Name: "John",
						Age:  42,
					},
				},
				Custom2: []*CustomType{
					{
						Name: "Jane",
						Age:  21,
					},
				},
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "map type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["custom"] = []string{`{"name":"John","age":42}`}
				return nil
			},
			v: &struct {
				Custom map[string]any `form:"custom"`
			}{},
			wantValue: &struct {
				Custom map[string]any `form:"custom"`
			}{
				Custom: map[string]any{
					"name": "John",
					"age":  float64(42),
				},
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "skip unexported fields",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["name"] = []string{`John`}
				return nil
			},
			v: &struct {
				name string `form:"name"`
			}{},
			wantValue: &struct {
				name string `form:"name"`
			}{
				name: "",
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "skip ignored fields",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["Name"] = []string{`John`}
				return nil
			},
			v: &struct {
				Name string `form:"-"`
			}{},
			wantValue: &struct {
				Name string `form:"-"`
			}{
				Name: "",
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "skip not existing file",
			fillMulipartForm: func(mpf *multipart.Form) error {
				return nil
			},
			v: &struct {
				File *multipart.FileHeader `form:"file"`
			}{},
			wantValue: &struct {
				File *multipart.FileHeader `form:"file"`
			}{
				File: nil,
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "skip not existing value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				return nil
			},
			v: &struct {
				Name string `form:"name"`
			}{},
			wantValue: &struct {
				Name string `form:"name"`
			}{
				Name: "",
			},
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when non-pointer type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				return nil
			},
			v:         struct{}{},
			wantValue: struct{}{},
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrValueMustBePointer) {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when non-struct type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				return nil
			},
			v:         ptr(""),
			wantValue: ptr(""),
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrValueMustBeStruct) {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when nil value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				return nil
			},
			v:         (*CustomType)(nil),
			wantValue: (*CustomType)(nil),
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrValueCannotBeNil) {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when invalid field type",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["invalid"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Invalid func() `form:"invalid"`
			}{},
			wantValue: &struct {
				Invalid func() `form:"invalid"`
			}{},
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrInvalidFieldType) {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when invalid json value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["invalid"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Invalid struct{} `form:"invalid"`
			}{},
			wantValue: &struct {
				Invalid struct{} `form:"invalid"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Invalid" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing int value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value int `form:"value"`
			}{},
			wantValue: &struct {
				Value int `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing uint value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value uint `form:"value"`
			}{},
			wantValue: &struct {
				Value uint `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing float value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value float64 `form:"value"`
			}{},
			wantValue: &struct {
				Value float64 `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing bool value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value bool `form:"value"`
			}{},
			wantValue: &struct {
				Value bool `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing complex value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value complex64 `form:"value"`
			}{},
			wantValue: &struct {
				Value complex64 `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing pointer value",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value *complex64 `form:"value"`
			}{},
			wantValue: &struct {
				Value *complex64 `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
		{
			name: "error when parsing custom type implements encoding.TextUnmarshaler",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["value"] = []string{"invalid"}
				return nil
			},
			v: &struct {
				Value failedCustomText `form:"value"`
			}{},
			wantValue: &struct {
				Value failedCustomText `form:"value"`
			}{},
			checkError: func(t *testing.T, err error) {
				var terr ErrParseFailed
				if !errors.As(err, &terr) || terr.Field != "Value" {
					t.Fatal("unexpected error:", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mpf := &multipart.Form{
				Value: make(map[string][]string),
				File:  make(map[string][]*multipart.FileHeader),
			}
			err := tt.fillMulipartForm(mpf)
			if err != nil {
				t.Fatal(err)
			}
			err = Convert(mpf, tt.v)
			tt.checkError(t, err)
			if !reflect.DeepEqual(tt.v, tt.wantValue) {
				t.Errorf("Convert() got = %v, want %v", tt.v, tt.wantValue)
				t.Fail()
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
