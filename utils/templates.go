package utils

import "html/template"

// ParseTemplates parses all the templates specified by the list of files
func ParseTemplates(templatePaths []string) (*template.Template, error) {
	templates, err := template.ParseFiles(templatePaths...)
	if err != nil {
		return nil, err
	}

	return templates, nil
}
