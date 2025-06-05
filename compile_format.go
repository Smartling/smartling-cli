package main

import (
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/reconquest/hierr-go"
)

var (
	// because we can have specific formats for different file types defined
	// in config file, we need to cache templates to prevent overhead in
	// runtime
	compiledFormatsCache = struct {
		sync.Mutex

		contents map[string]*Format
	}{
		contents: map[string]*Format{},
	}
)

func compileFormat(definition string) (*Format, error) {
	compiledFormatsCache.Lock()
	defer compiledFormatsCache.Unlock()

	if format, ok := compiledFormatsCache.contents[definition]; ok {
		return format, nil
	}

	definition = strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
	).Replace(definition)

	funcs := template.FuncMap{
		"name": func(path string) string {
			return strings.TrimSuffix(path, filepath.Ext(path))
		},

		"ext": func(path string) string {
			return filepath.Ext(path)
		},
	}

	var (
		format Format
		err    error
	)

	format.Source = definition
	format.Template, err = template.New("format").Funcs(funcs).Option(
		"missingkey=error",
	).Parse(
		definition,
	)
	if err != nil {
		return nil, clierror.NewError(
			hierr.Errorf(
				err,
				"failed to compile format template",
			),

			"Check template syntax accordingly to text/template "+
				"documentation:\n\thttps://golang.org/pkg/text/template/",
		)
	}

	compiledFormatsCache.contents[definition] = &format

	return &format, nil
}
