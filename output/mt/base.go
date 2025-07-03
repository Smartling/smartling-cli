package mt

import "github.com/charmbracelet/bubbles/table"

type Model struct {
	Headers        []table.Column
	Data           []table.Row
	RowByHeader    RowByHeaderName
	OutputFormat   string
	OutputTemplate string
}
