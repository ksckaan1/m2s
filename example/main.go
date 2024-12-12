package main

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/ksckaan1/m2s"
)

func main() {
	http.HandleFunc("POST /", rootPage)

	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type ReqBody struct {
	Name    string                `form:"name"`
	Age     int                   `form:"age"`
	Hobbies []string              `form:"hobbies"`
	File    *multipart.FileHeader `form:"file"`
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mpf := r.MultipartForm

	reqBody := &ReqBody{}

	err = m2s.Convert(mpf, reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	fmt.Printf("reqBody: %+v\n", reqBody)
}

type Supported struct {
	Text     string     `form:"text"`     // string | *string
	Decimal  int        `form:"number"`   // all int | int8 | int16 | int32 | int64
	Floating float64    `form:"floating"` // all float types
	Complex  complex128 `form:"complex"`
	Bool     bool       `form:"bool"`
}
