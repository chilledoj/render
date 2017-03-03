package render

import (
	"fmt"
	"html/template"
	"io"
	"sync"

	"github.com/oxtoacart/bpool"
)

// TemplateManager holds all of the templates and provides the buffered render function
type TemplateManager struct {
	bufpool *bpool.BufferPool
	tmpls   map[string]*template.Template
	mu      sync.Mutex
}

// NewTM creates a new TemplateManager with bufferpool initialised
func NewTM() *TemplateManager {
	tm := TemplateManager{
		bufpool: bpool.NewBufferPool(64),
		tmpls:   make(map[string]*template.Template),
	}
	return &tm
}

// AddTemplate adds a *template.Template to the internal map using name as key
func (tm *TemplateManager) AddTemplate(name string, tmpl *template.Template) error {
	tm.mu.Lock()
	tm.tmpls[name] = tmpl
	tm.mu.Unlock()
	return nil
}

// Render writes the template defined by name out to the io.Writer using data
func (tm *TemplateManager) Render(w io.Writer, name string, data interface{}) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	//log.Println(templates)
	tmpl, ok := tm.tmpls[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	buf := tm.bufpool.Get()
	defer tm.bufpool.Put(buf)

	err := tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	buf.WriteTo(w)
	return nil
}
