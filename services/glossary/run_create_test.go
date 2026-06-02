package glossary

import (
	"errors"
	"testing"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

func TestCreateParams_Validate(t *testing.T) {
	validAccount := uid.AccountUID("test-account-uid")

	tests := []struct {
		name    string
		params  CreateParams
		wantErr bool
	}{
		{
			name: "valid minimal",
			params: CreateParams{
				AccountUID:   validAccount,
				GlossaryName: "My Glossary",
				LocaleIDs:    []string{"en-US"},
			},
			wantErr: false,
		},
		{
			name: "valid with fallback locales",
			params: CreateParams{
				AccountUID:   validAccount,
				GlossaryName: "My Glossary",
				LocaleIDs:    []string{"en-US", "es-ES"},
				FallbackLocales: []FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES", "es-MX"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty account UID",
			params:  CreateParams{GlossaryName: "My Glossary", LocaleIDs: []string{"en-US"}},
			wantErr: true,
		},
		{
			name:    "empty glossary name",
			params:  CreateParams{AccountUID: validAccount, LocaleIDs: []string{"en-US"}},
			wantErr: true,
		},
		{
			name:    "empty locale IDs",
			params:  CreateParams{AccountUID: validAccount, GlossaryName: "My Glossary"},
			wantErr: true,
		},
		{
			name: "fallback locale missing FallbackLocaleID",
			params: CreateParams{
				AccountUID:   validAccount,
				GlossaryName: "My Glossary",
				LocaleIDs:    []string{"en-US"},
				FallbackLocales: []FallbackLocale{
					{FallbackLocaleID: "", LocaleIDs: []string{"es-ES"}},
				},
			},
			wantErr: true,
		},
		{
			name: "fallback locale missing LocaleIDs",
			params: CreateParams{
				AccountUID:   validAccount,
				GlossaryName: "My Glossary",
				LocaleIDs:    []string{"en-US"},
				FallbackLocales: []FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: nil},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.params.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_RunCreate(t *testing.T) {
	ctx := t.Context()

	validParams := CreateParams{
		AccountUID:   uid.AccountUID("test-account-uid"),
		GlossaryName: "My Glossary",
		LocaleIDs:    []string{"en-US", "es-ES"},
	}
	apiResponse := glossaryapi.CreateGlossaryResponse{
		Code:         200,
		GlossaryUID:  "test-glossary-uid",
		AccountUID:   string(validParams.AccountUID),
		GlossaryName: validParams.GlossaryName,
	}

	tests := []struct {
		name    string
		setup   func(*sdkmocks.MockGlossary)
		params  CreateParams
		wantErr bool
		check   func(*testing.T, CreateOutput)
	}{
		{
			name: "success",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().Create(ctx, validParams.AccountUID, toAPICreateParams(validParams)).
					Return(apiResponse, nil)
			},
			params: validParams,
			check: func(t *testing.T, got CreateOutput) {
				if got.GlossaryUID != apiResponse.GlossaryUID {
					t.Errorf("GlossaryUID = %v, want %v", got.GlossaryUID, apiResponse.GlossaryUID)
				}
				if got.GlossaryName != apiResponse.GlossaryName {
					t.Errorf("GlossaryName = %v, want %v", got.GlossaryName, apiResponse.GlossaryName)
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name:    "validation error — no locale IDs",
			setup:   func(m *sdkmocks.MockGlossary) {},
			params:  CreateParams{AccountUID: uid.AccountUID("test-account-uid"), GlossaryName: "My Glossary"},
			wantErr: true,
		},
		{
			name: "API error",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().Create(ctx, validParams.AccountUID, toAPICreateParams(validParams)).
					Return(glossaryapi.CreateGlossaryResponse{}, errors.New("API error"))
			},
			params:  validParams,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m)
			got, err := service{glossaryApi: m}.RunCreate(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
