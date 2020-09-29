package handfile

import (
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	path := "filePage/" + filename
	return ioutil.WriteFile(path, p.Body, 0600)
}
func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("filePage/" + filename)
	return &Page{Title: title, Body: body}, err
}
func AddPage(page *Page) error {
	p1 := &Page{Title: page.Title, Body: []byte(page.Body)}
	err := p1.save()
	return err
}
