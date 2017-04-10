package render_test

import (
	"bytes"
	"html/template"
	"os"
	"sync"
	"testing"

	"github.com/chilledoj/render"
)

func TestManager(t *testing.T) {
	tm := render.NewTM()
	if tm == nil {
		t.Fatal("render.NewTM does not return an instance of TemplateManager")
	}
	validTemplate, err := template.New("test").Parse(`<one>{{.Test}}</one>`)
	if err != nil {
		t.Fatalf("Invalid Test template", err)
	}
	if err = tm.AddTemplate("test", validTemplate); err != nil {
		t.Errorf("Template can not be added to TemplateManager map", err)
	}
	data := struct {
		Test string
	}{"TEST"}
	b := new(bytes.Buffer)

	if err = tm.Render(b, "test", data); err != nil {
		t.Fatalf("Render returned an error ", err, b)
	}
	returnedString := b.String()
	if returnedString != "<one>TEST</one>" {
		t.Fatalf("Invalid returned string (%s)", returnedString)
	}
	if err = tm.Render(b, "NOTAVAILABLE", data); err == nil {
		t.Fatal("Invalid template does not result in error")
	}

	invalidData := "TEST"
	if err = tm.Render(b, "test", invalidData); err == nil {
		t.Fatal("No error on supplying invalid data")
	}
}

// ExampleTemplateManager
func ExampleTemplateManager() {
	tm := render.NewTM()
	validTemplate, err := template.New("").Parse(`<one>{{.Test}}</one>`)
	if err != nil {
		panic(err)
	}
	if err = tm.AddTemplate("one", validTemplate); err != nil {
		panic(err)
	}
	data := struct {
		Test string
	}{"Template One"}

	if err = tm.Render(os.Stdout, "one", data); err != nil {
		panic(err)
	}
	// Output: <one>Template One</one>
}

func BenchmarkSequential(b *testing.B) {
	tm := render.NewTM()
	validTemplate, _ := template.New("").Parse(`<one>{{.Test}}</one>`)
	if err := tm.AddTemplate("test", validTemplate); err != nil {
		b.Errorf("Template can not be added to TemplateManager map", err)
	}
	data := struct {
		Test string
	}{"TEST"}
	buf := new(bytes.Buffer)

	for i := 0; i < b.N; i++ {
		if err := tm.Render(buf, "test", data); err != nil {
			b.Fatalf("Render returned an error ", err, b)
		}
	}
}

func BenchmarkGoRoutine_AllocBuf(b *testing.B) {
	tm := render.NewTM()
	validTemplate, _ := template.New("").Parse(`<one>{{.Test}}</one>`)
	if err := tm.AddTemplate("test", validTemplate); err != nil {
		b.Errorf("Template can not be added to TemplateManager map", err)
	}
	data := struct {
		Test string
	}{"TEST"}

	for i := 0; i < b.N; i++ {
		go func() {
			buf := new(bytes.Buffer)
			if err := tm.Render(buf, "test", data); err != nil {
				b.Fatalf("Render returned an error ", err, b)
			}
		}()
	}
}

func BenchmarkGoRoutine_SingleBuf(b *testing.B) {
	tm := render.NewTM()
	validTemplate, _ := template.New("").Parse(`<one>{{.Test}}</one>`)
	if err := tm.AddTemplate("test", validTemplate); err != nil {
		b.Errorf("Template can not be added to TemplateManager map", err)
	}
	data := struct {
		Test string
	}{"TEST"}

	buf := new(bytes.Buffer)
	var lock sync.Mutex
	for i := 0; i < b.N; i++ {
		go func() {
			lock.Lock()
			if err := tm.Render(buf, "test", data); err != nil {
				b.Fatalf("Render returned an error ", err, b)
			}
			lock.Unlock()
		}()
	}
}
