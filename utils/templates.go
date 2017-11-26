package utils

import "html/template"

// ParseTemplates parses all the templates specified by the list of files
func ParseTemplates(templatePaths []string) (templates *template.Template, err error) {
	templates, err = template.ParseFiles(templatePaths...)
	return
}
