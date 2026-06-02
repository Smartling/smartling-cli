package glossaryresolver

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// glossaryUIDPattern matches a Smartling glossary UID — a canonical
// 8-4-4-4-12 hex UUID.
var glossaryUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// GetGlossaryUID resolves either a glossary UID or a glossary name into a UID.
// When the input matches the UID shape the API is asked for that glossary
// directly; on miss (or for non-UID-shaped input) the account is searched by
// name and the first matching glossary's UID is returned.
func GetGlossaryUID(ctx context.Context, api glossaryapi.Glossary, accountUID uid.AccountUID, glossaryUIDOrName string) (string, error) {
	if strings.TrimSpace(glossaryUIDOrName) == "" {
		return "", glossaryapi.ErrGlossaryNotFound
	}
	if glossaryUIDPattern.MatchString(glossaryUIDOrName) {
		gl, err := api.Get(ctx, accountUID, glossaryUIDOrName)
		switch {
		case err == nil:
			if strings.TrimSpace(gl.GlossaryUID) == "" {
				return "", glossaryapi.ErrGlossaryNotFound
			}
			return gl.GlossaryUID, nil
		case errors.Is(err, glossaryapi.ErrGlossaryNotFound):
			// 12-char input wasn't a UID — could still be a glossary name, fall through
		default:
			return "", fmt.Errorf("get glossary by UID %q: %w", glossaryUIDOrName, err)
		}
	}

	glossaries, err := api.GetByName(ctx, accountUID, glossaryUIDOrName)
	if err != nil {
		return "", fmt.Errorf("search glossaries by name %q: %w", glossaryUIDOrName, err)
	}
	if len(glossaries) == 0 {
		return "", glossaryapi.ErrGlossaryNotFound
	}
	for _, glossary := range glossaries {
		if glossary.Name == glossaryUIDOrName && strings.TrimSpace(glossary.GlossaryUID) != "" {
			return glossary.GlossaryUID, nil
		}
	}
	first := glossaries[0]
	if strings.TrimSpace(first.GlossaryUID) == "" {
		return "", glossaryapi.ErrGlossaryNotFound
	}
	return first.GlossaryUID, nil
}
