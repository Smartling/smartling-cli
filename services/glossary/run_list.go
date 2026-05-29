package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// ListParams carries the glossary-list request from CLI to service.
// Name is an optional filter passed to the API as `glossaryName`; an empty
// value lists all glossaries under the account.
type ListParams struct {
	AccountUID uid.AccountUID
	Name       string
}

// Validate checks that ListParams carry the fields required by the API.
func (p ListParams) Validate() error {
	return p.AccountUID.Validate()
}

// GlossaryItem is a single row in a glossary list.
type GlossaryItem struct {
	GlossaryUID string
	Name        string
	Description string
	LocaleIDs   []string
}

// ListOutput represents the result of a glossary list.
type ListOutput struct {
	Glossaries []GlossaryItem
	JSON       []byte
}

// JSONBytes returns the raw JSON payload of the list response.
func (o ListOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary of the list.
func (o ListOutput) SimpleLines() []string {
	if len(o.Glossaries) == 0 {
		return []string{"No glossaries found."}
	}
	lines := make([]string, 0, len(o.Glossaries))
	for _, g := range o.Glossaries {
		lines = append(lines, fmt.Sprintf("%s  %s  %s", g.GlossaryUID, g.Name, g.Description))
	}
	return lines
}

// TableData returns the list as one column per field, one row per glossary.
func (o ListOutput) TableData() ([]string, [][]string) {
	headers := []string{"GLOSSARY UID", "NAME", "DESCRIPTION", "LOCALES"}
	rows := make([][]string, 0, len(o.Glossaries))
	for _, g := range o.Glossaries {
		rows = append(rows, []string{
			g.GlossaryUID,
			g.Name,
			g.Description,
			strings.Join(g.LocaleIDs, ","),
		})
	}
	return headers, rows
}

// RunList lists glossaries under the account, optionally filtered by name.
func (s service) RunList(ctx context.Context, params ListParams) (ListOutput, error) {
	if err := params.Validate(); err != nil {
		return ListOutput{}, fmt.Errorf("invalid list params: %w", err)
	}

	resp, err := s.glossaryApi.GetByName(ctx, string(params.AccountUID), params.Name)
	if err != nil {
		return ListOutput{}, fmt.Errorf("failed to list glossaries: %w", err)
	}

	return toListOutput(resp), nil
}

func toListOutput(resp []api.ReadGlossaryResponse) ListOutput {
	items := make([]GlossaryItem, 0, len(resp))
	for _, g := range resp {
		items = append(items, GlossaryItem{
			GlossaryUID: g.GlossaryUid,
			Name:        g.Name,
			Description: g.Description,
			LocaleIDs:   g.LocaleIDs,
		})
	}

	out := ListOutput{Glossaries: items}
	if b, err := json.Marshal(items); err == nil {
		out.JSON = b
	}
	return out
}
