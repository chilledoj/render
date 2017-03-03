package render_test

import (
	"bytes"
	"html/template"
	"os"
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
	err = tm.AddTemplate("test", validTemplate)
	if err != nil {
		t.Errorf("Template can not be added to TemplateManager map", err)
	}
	data := struct {
		Test string
	}{"TEST"}
	b := new(bytes.Buffer)

	err = tm.Render(b, "test", data)
	if err != nil {
		t.Fatalf("Render returned an error ", err, b)
	}
	returnedString := b.String()
	if returnedString != "<one>TEST</one>" {
		t.Fatalf("Invalid returned string (%s)", returnedString)
	}
	err = tm.Render(b, "NOTAVAILABLE", data)
	if err == nil {
		t.Fatal("Invalid template does not result in error")
	}

	invalidData := "TEST"
	err = tm.Render(b, "test", invalidData)
	if err == nil {
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
	err = tm.AddTemplate("one", validTemplate)
	if err != nil {
		panic(err)
	}
	data := struct {
		Test string
	}{"Template One"}

	err = tm.Render(os.Stdout, "one", data)
	if err != nil {
		panic(err)
	}
	// Output: <one>Template One</one>
}
