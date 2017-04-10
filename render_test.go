package render

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
			if err := BufferedRender(tt.tmpl, b, tt.templateName, tt.data); err != nil && !tt.wantErr {
				t.Errorf("Unexpected error with %s: %v", tt.templateName, err)
			}
		})
	}
}

func ExampleBufferedRender() {
	tmpl, err := template.New("one").Parse(`<one>{{.Test}}</one>`)
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	if err := BufferedRender(tmpl, b, "two", nil); err != nil {
		fmt.Printf("Error rendering: %s", err)
		return
	}
	b.WriteTo(os.Stdout)
	// outputs: Error rendering: html/template: "two" is undefined
}
func ExampleBufRenderHandler() {
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
	// outputs: <html><body><h1>Hello World</h1><p>Testing Templates</p></body></html>
}
func hello() http.HandlerFunc {
	tmpl, err := template.New("hello").Parse(`<html><body><h1>Hello World</h1><p>{{.Test}}</p></body></html>`)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := BufferedRender(tmpl, w, "hello", struct{ Test string }{"Testing Templates"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
