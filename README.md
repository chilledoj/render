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
or in the context of an http.HandlerFunc
```go
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
```