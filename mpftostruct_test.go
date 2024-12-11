package mpftostruct

import (
	"errors"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name             string
		fillMulipartForm func(*multipart.Form) error
		v                any
		wantValue        any
		wantErr          error
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
			wantErr: nil,
		},
		{
			name: "primitive types list",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["key1"] = []string{"value1", "value2"}
				mpf.Value["key2"] = []string{"42", "43"}
				mpf.Value["key3"] = []string{"42", "43"}
				mpf.Value["key4"] = []string{"42.5", "43.5"}
				mpf.Value["key5"] = []string{"10+11i", "10+12i"}
				mpf.Value["key6"] = []string{"true", "false"}
				return nil
			},
			v: &struct {
				Key1 []string    `form:"key1"`
				Key2 []int32     `form:"key2"`
				Key3 []uint16    `form:"key3"`
				Key4 []float32   `form:"key4"`
				Key5 []complex64 `form:"key5"`
				Key6 []bool      `form:"key6"`
			}{},
			wantValue: &struct {
				Key1 []string    `form:"key1"`
				Key2 []int32     `form:"key2"`
				Key3 []uint16    `form:"key3"`
				Key4 []float32   `form:"key4"`
				Key5 []complex64 `form:"key5"`
				Key6 []bool      `form:"key6"`
			}{
				Key1: []string{"value1", "value2"},
				Key2: []int32{42, 43},
				Key3: []uint16{42, 43},
				Key4: []float32{42.5, 43.5},
				Key5: []complex64{10 + 11i, 10 + 12i},
				Key6: []bool{true, false},
			},
			wantErr: nil,
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
			wantErr: nil,
		},
		{
			name: "pointer types list",
			fillMulipartForm: func(mpf *multipart.Form) error {
				mpf.Value["key1"] = []string{"value1", "value2"}
				mpf.Value["key2"] = []string{"42", "43"}
				mpf.Value["key3"] = []string{"42", "43"}
				mpf.Value["key4"] = []string{"42.5", "43.5"}
				mpf.Value["key5"] = []string{"10+11i", "10+12i"}
				mpf.Value["key6"] = []string{"true", "false"}
				return nil
			},
			v: &struct {
				Key1 []*string    `form:"key1"`
				Key2 []*int32     `form:"key2"`
				Key3 []*uint16    `form:"key3"`
				Key4 []*float32   `form:"key4"`
				Key5 []*complex64 `form:"key5"`
				Key6 []*bool      `form:"key6"`
			}{},
			wantValue: &struct {
				Key1 []*string    `form:"key1"`
				Key2 []*int32     `form:"key2"`
				Key3 []*uint16    `form:"key3"`
				Key4 []*float32   `form:"key4"`
				Key5 []*complex64 `form:"key5"`
				Key6 []*bool      `form:"key6"`
			}{
				Key1: []*string{ptr("value1"), ptr("value2")},
				Key2: []*int32{ptr[int32](42), ptr[int32](43)},
				Key3: []*uint16{ptr[uint16](42), ptr[uint16](43)},
				Key4: []*float32{ptr[float32](42.5), ptr[float32](43.5)},
				Key5: []*complex64{ptr[complex64](10 + 11i), ptr[complex64](10 + 12i)},
				Key6: []*bool{ptr(true), ptr(false)},
			},
			wantErr: nil,
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
			wantErr: nil,
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
			wantErr: nil,
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
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				t.Fail()
			}
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
