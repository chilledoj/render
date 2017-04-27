# Render
A renderer for Golang projects which need to handle errors properly.

```go
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
if err := render.BufferedRender(tmpl, b, "one", nil); err != nil {
	fmt.Printf("Error rendering: %s", err)
	return
}
b.WriteTo(os.Stdout)
// Output: <one></one>
```
or in the context of an http.HandlerFunc
```go
func main() {
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
```