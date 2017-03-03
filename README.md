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
