package mt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"text/template"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	outtable "github.com/charmbracelet/lipgloss/table"
)

// OutputFormat defines behaviour to format and render data
type OutputFormat interface {
	FormatAndRender(headers []table.Column, data []table.Row)
}

// GetOutputFormat returns OutputFormat for given string
func GetOutputFormat(outputFormat, outputTemplate string) OutputFormat {
	switch outputFormat {
	case "table":
		return TableOutputFormat{}
	case "json":
		return JsonOutputFormat{}
	case "simple":
		return SimpleOutputFormat{template: outputTemplate}
	}
	return SimpleOutputFormat{template: outputTemplate}
}

// TableOutputFormat is table output format for rendering data as a formatted table
type TableOutputFormat struct{}

// FormatAndRender renders the data as a table
func (t TableOutputFormat) FormatAndRender(headers []table.Column, data []table.Row) {
	hh := make([]string, len(headers))
	for i, h := range headers {
		hh[i] = h.Title
	}
	tbl := outtable.New().
		Border(lipgloss.NormalBorder()).
		Headers(hh...)
	for _, row := range data {
		tbl.Row(row...)
	}
	fmt.Println(tbl)
}

// JsonOutputFormat is json output format for rendering data as json
type JsonOutputFormat struct{}

// FormatAndRender marshals table data into JSON and prints it.
func (j JsonOutputFormat) FormatAndRender(headers []table.Column, data []table.Row) {
	hh := make([]string, len(headers))
	for i, h := range headers {
		hh[i] = h.Title
	}
	var result []map[string]string
	for _, row := range data {
		m := make(map[string]string)
		for i, val := range row {
			if i < len(hh) {
				m[hh[i]] = val
			}
		}
		result = append(result, m)
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}

// SimpleOutputFormat is simple output format for rendering data using a text template
type SimpleOutputFormat struct {
	template string
}

// FormatAndRender formats the data rows using the stored template string
// and outputs the rendered result
func (s SimpleOutputFormat) FormatAndRender(headers []table.Column, data []table.Row) {
	funcMap := template.FuncMap{
		"name": func(f string) string {
			return filepath.Base(f[:len(f)-len(filepath.Ext(f))])
		},
		"ext": func(f string) string {
			return filepath.Ext(f)
		},
		"dir": func(f string) string {
			return filepath.Dir(f)
		},
	}
	tmpl, err := template.New("rowformat").Funcs(funcMap).Parse(s.template)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range data {
		rowMap := make(map[string]string)
		for i, col := range headers {
			rowMap[col.Title] = row[i]
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, rowMap); err != nil {
			log.Fatal(err)
		}

		fmt.Print(buf.String())
	}
}
