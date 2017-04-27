package render_test

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/chilledoj/render"
)

func TestBufferedRender(t *testing.T) {
	tmpl, err := template.New("one").Parse(`<one>{{.Test}}</one>`)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testname     string
		tmpl         *template.Template
		templateName string
		data         interface{}
		wantErr      bool
	}{
		{"Nil Data", tmpl, "one", nil, false},
		{"Valid Data", tmpl, "one", struct{ Test string }{"Testing"}, false},
		{"Invalid template name", tmpl, "two", nil, true},
	}
	var bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t2 *testing.T) {
			b := bufPool.Get().(*bytes.Buffer)
			defer bufPool.Put(b)
			if err := render.BufferedRender(tt.tmpl, b, tt.templateName, tt.data); err != nil && !tt.wantErr {
				t.Errorf("Unexpected error with %s: %v", tt.templateName, err)
			}
		})
	}
}

// ExampleRender_BufferedRender
func ExampleRender_BufferedRender() {
	tmpl, err := template.New("one").Parse(`<one>{{.Test}}</one>`)
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	if err := render.BufferedRender(tmpl, b, "one", nil); err != nil {
		fmt.Printf("Error rendering: %s", err)
		return
	}
	b.WriteTo(os.Stdout)
	// Output: <one></one>
}
func ExampleRender_BufferedRender_Handler() {
	ts := httptest.NewServer(hello())
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
	// Output: <html><body><h1>Hello World</h1><p>Testing Templates</p></body></html>
}
func hello() http.HandlerFunc {
	tmpl, err := template.New("hello").Parse(`<html><body><h1>Hello World</h1><p>{{.Test}}</p></body></html>`)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := render.BufferedRender(tmpl, w, "hello", struct{ Test string }{"Testing Templates"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
