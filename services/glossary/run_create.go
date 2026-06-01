package glossary

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/glossary"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// CreateParams defines full glossary-create request from CLI to service.
// Field semantics match the Smartling Glossary Create API
// (https://api-reference.smartling.com/#tag/Glossary-API/operation/createGlossary).
type CreateParams struct {
	AccountUID       uid.AccountUID
	GlossaryName     string
	Description      string
	VerificationMode bool
	LocaleIDs        []string
	FallbackLocales  []FallbackLocale
}

// FallbackLocale maps a fallback source locale to the set of target locales
// that should fall back to it.
type FallbackLocale struct {
	FallbackLocaleID string
	LocaleIDs        []string
}

// Validate checks that CreateParams carry the fields required by the API.
func (p CreateParams) Validate() error {
	if err := p.AccountUID.Validate(); err != nil {
		return err
	}
	if p.GlossaryName == "" {
		return smerror.ErrEmptyParam("GlossaryName")
	}
	if len(p.LocaleIDs) == 0 {
		return smerror.ErrEmptyParam("LocaleIDs")
	}
	for _, fallbackLocale := range p.FallbackLocales {
		if fallbackLocale.FallbackLocaleID == "" {
			return smerror.ErrEmptyParam("FallbackLocale.FallbackLocaleID")
		}
		if len(fallbackLocale.LocaleIDs) == 0 {
			return smerror.ErrEmptyParam("FallbackLocale.LocaleIDs")
		}
	}
	return nil
}

// CreateOutput represents the result of a glossary create.
type CreateOutput struct {
	Code         int
	GlossaryUID  string
	AccountUID   string
	GlossaryName string
	JSON         []byte
}

// JSONBytes returns the raw JSON payload of the create response.
func (o CreateOutput) JSONBytes() []byte { return o.JSON }

// SimpleLines returns a human-readable summary of the create.
func (o CreateOutput) SimpleLines() []string {
	return []string{
		fmt.Sprintf("Glossary UID:  %s", o.GlossaryUID),
		fmt.Sprintf("Account UID:   %s", o.AccountUID),
		fmt.Sprintf("Glossary name: %s", o.GlossaryName),
	}
}

// TableData returns the create summary with one column per field and a
// single row of values.
func (o CreateOutput) TableData() ([]string, [][]string) {
	headers := []string{"GLOSSARY UID", "ACCOUNT UID", "GLOSSARY NAME"}
	rows := [][]string{
		{o.GlossaryUID, o.AccountUID, o.GlossaryName},
	}
	return headers, rows
}

func (s service) RunCreate(ctx context.Context, params CreateParams) (CreateOutput, error) {
	if err := params.Validate(); err != nil {
		return CreateOutput{}, fmt.Errorf("failed parameters validation: %w", err)
	}
	apiParams := toAPICreateParams(params)
	resp, err := s.glossaryApi.Create(ctx, params.AccountUID, apiParams)
	if err != nil {
		return CreateOutput{}, fmt.Errorf("glossary create API call failed: %w", err)
	}
	return toCreateOutput(resp), nil
}

func toCreateOutput(resp api.CreateGlossaryResponse) CreateOutput {
	res := CreateOutput{
		Code:         resp.Code,
		GlossaryUID:  resp.GlossaryUID,
		AccountUID:   resp.AccountUID,
		GlossaryName: resp.GlossaryName,
	}

	summary := struct {
		GlossaryUID  string `json:"glossaryUid"`
		AccountUID   string `json:"accountUid"`
		GlossaryName string `json:"glossaryName"`
	}{
		GlossaryUID:  resp.GlossaryUID,
		AccountUID:   resp.AccountUID,
		GlossaryName: resp.GlossaryName,
	}
	b, err := json.Marshal(summary)
	if err != nil {
		rlog.Errorf("failed to marshal create output to JSON: %v", err)
		return res
	}
	res.JSON = b
	return res
}

func toAPICreateParams(p CreateParams) api.CreateGlossaryRequest {
	fallbackLocales := make([]api.FallbackLocale, 0, len(p.FallbackLocales))
	for _, fallbackLocale := range p.FallbackLocales {
		fallbackLocales = append(fallbackLocales, api.FallbackLocale{
			FallbackLocaleID: fallbackLocale.FallbackLocaleID,
			LocaleIDs:        fallbackLocale.LocaleIDs,
		})
	}
	return api.CreateGlossaryRequest{
		GlossaryName:     p.GlossaryName,
		Description:      p.Description,
		VerificationMode: p.VerificationMode,
		LocaleIDs:        p.LocaleIDs,
		FallbackLocales:  fallbackLocales,
	}
}
