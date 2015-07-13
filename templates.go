package main

import "html/template"

var indexTemplate = initializeTemplate("templates/index.html")

type indexTemplateVars struct {
	Name string
}

var errorTemplate = initializeTemplate("templates/error.html")

type errorTemplateVars struct {
	Code    int
	Message string
}

func initializeTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles("templates/layout.html", file))
}
