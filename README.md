# Render
A renderer for Golang projects which need to handle errors properly.

```
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
```
## Simpler - just buffered renderer
There is also the following function provided for a simpler straight forward buffered template executor.
```go
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
```