/*
Package render is a buffered renderer for Golang html/template
*/
package render

import (
	"bytes"
	"html/template"
	"io"
	"sync"
)

var globalBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// BufferedRender is a shorthand way of buffering render to
func BufferedRender(tmpl *template.Template, w io.Writer, name string, data interface{}) error {
	b := globalBufPool.Get().(*bytes.Buffer)
	defer globalBufPool.Put(b)
	b.Reset()

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		return err
	}
	b.WriteTo(w)
	return nil
}
