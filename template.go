package main

import (
	"io"
	"log"
	"mime"
	"path/filepath"
	"strings"
	"text/template"
)

// Template is a resource rendered from a go template.
type Template struct {
	template     *template.Template
	templateName string
	variables    *VariableMap
}

func newTemplate(fileName string) *Template {
	name := filepath.Base(fileName)
	template := template.Must(
		template.New(name).
			Funcs(template.FuncMap{"lower": strings.ToLower, "shortHash": ShortHash}).
			ParseFiles(fileName),
	)

	return &Template{template: template, templateName: name}
}

// MimeType derives the mime type of the resource from the template file
// extension. An empty string is returned if there is no extension or the
// mime type for the extension is not known.
func (r *Template) MimeType() string {
	ext := filepath.Ext(r.templateName)

	if ext == "" {
		return ""
	}

	return mime.TypeByExtension(ext)
}

// Render renders the template.
func (r *Template) Render(out io.Writer, params ParamMap, vars VariableMap) error {
	log.Printf("Render template: %v %v", r.templateName, params)
	return r.template.ExecuteTemplate(out, r.templateName, struct{
		Params ParamMap
		Vars VariableMap
	}{Params: params, Vars: vars})
}
