package utils

import (
	"html/template"
	"path/filepath"
)

// ParseTemplatesFromBinData parses all the templates specified by the list of files using go-bindata assets
func ParseTemplatesFromBinData(templateDir string) (*template.Template, error) {
	assets, err := AssetDir(templateDir)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(assets))
	for idx, asset := range assets {
		paths[idx] = "templates/" + asset
	}

	return parseTemplates(paths)
}

func parseTemplates(paths []string) (*template.Template, error) {
	if paths == nil || len(paths) == 0 {
		return nil, nil
	}
	t, err := parseTemplate(nil, paths[0])
	if err != nil {
		return nil, err
	}
	for _, p := range paths[1:] {
		_, err := parseTemplate(t, p)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func parseTemplate(t *template.Template, p string) (*template.Template, error) {
	templateName := filepath.Base(p)
	a, err := Asset(p)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return template.New(templateName).Parse(string(a))
	}
	_, err = t.New(templateName).Parse(string(a))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
