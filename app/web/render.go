package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	IsAuthenticated int
	ErrorData       []string
	Flash           string
	Warning         string
	Error           string
	API             string
	CSSVersion      string
}

//go:embed templates
var templateFS embed.FS

func (web *webapp) addDefaultData(td *templateData, r *http.Request) *templateData {
	return td
}

func (web *webapp) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	_, templateInMap := web.templateCache[templateToRender]

	if web.env == "production" && templateInMap {
		t = web.templateCache[templateToRender]
	} else {
		t, err = web.parseTemplate(page, templateToRender)
		if err != nil {
			fmt.Printf("Error parsing template: %v", err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = web.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		fmt.Printf("Error executing template: %v", err)
		return err
	}
	return nil
}

func (web *webapp) parseTemplate(page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)

	if err != nil {
		fmt.Printf("Error creating template: %v", err)
		return nil, err
	}

	web.templateCache[templateToRender] = t
	return t, nil
}
