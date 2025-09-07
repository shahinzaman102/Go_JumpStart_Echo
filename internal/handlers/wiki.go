package handlers

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
)

// Page represents a wiki page with a title and body
type Page struct {
	Title string
	Body  []byte // Page content as raw bytes
}

// save writes the Page's body to a file in data/<Title>.txt
func (p *Page) save() error {
	if err := os.MkdirAll("data", 0700); err != nil {
		return err
	}
	return os.WriteFile("data/"+p.Title+".txt", p.Body, 0600)
}

// loadPage reads a page from disk and returns a Page struct
func loadPage(title string) (*Page, error) {
	body, err := os.ReadFile("data/" + title + ".txt")
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// wikiTemplates is loaded lazily to avoid panics during init/tests
var wikiTemplates *template.Template

// LoadWikiTemplates parses the edit and view templates
func LoadWikiTemplates() error {
	var err error
	wikiTemplates, err = template.ParseFiles("templates/edit.html", "templates/view.html")
	return err
}

// validLink matches [[PageName]] style links in wiki text
var validLink = regexp.MustCompile(`\[(.+?)\]`)

// renderTemplate processes links and renders a wiki template
func renderTemplate(c echo.Context, tmpl string, p *Page) error {
	bodyStr := string(p.Body)

	// Convert [[PageName]] into <a href="/view/PageName">PageName</a>
	processed := validLink.ReplaceAllStringFunc(bodyStr, func(s string) string {
		m := validLink.FindStringSubmatch(s)
		link := m[1]
		escaped := url.PathEscape(link)
		return fmt.Sprintf(`<a href="/view/%s">%s</a>`, escaped, link)
	})

	return wikiTemplates.ExecuteTemplate(c.Response(), tmpl+".html", struct {
		Title string
		Body  template.HTML
	}{
		Title: p.Title,
		Body:  template.HTML(processed),
	})
}

// ViewWiki handles GET /view/:title
func ViewWiki(c echo.Context) error {
	title := c.Param("title")
	decodedTitle, _ := url.PathUnescape(title)

	p, err := loadPage(decodedTitle)
	if err != nil {
		return c.Redirect(302, "/edit/"+decodedTitle)
	}
	if err := renderTemplate(c, "view", p); err != nil {
		return c.String(500, err.Error())
	}
	return nil
}

// EditWiki handles GET /edit/:title
func EditWiki(c echo.Context) error {
	title := c.Param("title")
	decodedTitle, _ := url.PathUnescape(title)

	p, err := loadPage(decodedTitle)
	if err != nil {
		p = &Page{Title: decodedTitle} // new empty page
	}

	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	if err := tmpl.Execute(c.Response(), p); err != nil {
		return c.String(500, err.Error())
	}
	return nil
}

// SaveWiki handles POST /save/:title
func SaveWiki(c echo.Context) error {
	title := c.Param("title")
	decodedTitle, _ := url.PathUnescape(title)

	body := c.FormValue("body")
	p := &Page{Title: decodedTitle, Body: []byte(body)}

	if err := p.save(); err != nil {
		return c.String(500, err.Error())
	}

	return c.Redirect(302, "/view/"+decodedTitle)
}
