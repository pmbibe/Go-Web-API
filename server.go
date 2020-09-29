package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"./handfile"
)

var templates = template.Must(template.ParseFiles("templatePage/edit.html"))
var validPath = regexp.MustCompile("^/(editPage|addPage|view)/([a-zA-Z0-9]+)$")

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HomePage")
}
func handlerAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API")
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := handfile.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/editPage/"+title, http.StatusFound)
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func addPage(w http.ResponseWriter, r *http.Request, title string) {
	r.ParseForm()
	body := r.Form.Get("body")
	p1 := &handfile.Page{Title: title, Body: []byte(body)}
	err := handfile.AddPage(p1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *handfile.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func editPage(w http.ResponseWriter, r *http.Request, title string) {
	p, err := handfile.LoadPage(title)
	if err != nil {
		p = &handfile.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Print(m)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/api/", handlerAPI)
	mux.HandleFunc("/view/", makeHandler(viewHandler))
	mux.HandleFunc("/addPage/", makeHandler(addPage))
	mux.HandleFunc("/editPage/", makeHandler(editPage))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
