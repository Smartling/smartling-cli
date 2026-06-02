// Package static provides generic format strategies for rendering
// command results as JSON, human-readable text, or an ASCII table.
package static

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Renderable is the contract every payload must satisfy to be rendered
// by an OutputFormat. TableData returns column headers and one slice per
// row whose length matches len(headers).
type Renderable interface {
	JSONBytes() []byte
	SimpleLines() []string
	TableData() (headers []string, rows [][]string)
}

// OutputFormat formats and renders a payload of type T.
type OutputFormat[T Renderable] interface {
	FormatAndRender(data T)
}

// GetOutputFormat returns the OutputFormat selected by name. Unknown
// names fall back to the simple format.
func GetOutputFormat[T Renderable](name string) OutputFormat[T] {
	switch name {
	case "json":
		return JSONOutputFormat[T]{}
	case "table":
		return TableOutputFormat[T]{}
	default:
		return SimpleOutputFormat[T]{}
	}
}

// JSONOutputFormat prints the payload's JSON representation.
type JSONOutputFormat[T Renderable] struct{}

// FormatAndRender writes the JSON bytes from the payload to stdout.
func (JSONOutputFormat[T]) FormatAndRender(data T) {
	fmt.Println(string(data.JSONBytes()))
}

// SimpleOutputFormat prints each line returned by the payload.
type SimpleOutputFormat[T Renderable] struct{}

// FormatAndRender writes each simple line from the payload to stdout.
func (SimpleOutputFormat[T]) FormatAndRender(data T) {
	for _, line := range data.SimpleLines() {
		fmt.Println(line)
	}
}

// TableOutputFormat prints the payload as an ASCII table when the payload
// implements Tabular; otherwise it falls back to SimpleLines.
type TableOutputFormat[T Renderable] struct{}

// FormatAndRender writes the payload as an ASCII table.
func (TableOutputFormat[T]) FormatAndRender(data T) {
	headers, rows := data.TableData()
	tbl := table.New().
		Border(lipgloss.ASCIIBorder()).
		Headers(headers...)
	for _, row := range rows {
		tbl.Row(row...)
	}
	fmt.Println(tbl)
}
