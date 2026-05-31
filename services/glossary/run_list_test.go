package glossary

import (
	"errors"
	"testing"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

func TestListParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  ListParams
		wantErr bool
	}{
		{
			name:    "valid with no name filter",
			params:  ListParams{AccountUID: uid.AccountUID("abc123")},
			wantErr: false,
		},
		{
			name:    "valid with name filter",
			params:  ListParams{AccountUID: uid.AccountUID("abc123"), Name: "Marketing"},
			wantErr: false,
		},
		{
			name:    "empty account UID",
			params:  ListParams{Name: "Marketing"},
			wantErr: true,
		},
		{
			name:    "whitespace account UID",
			params:  ListParams{AccountUID: uid.AccountUID("   ")},
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

func Test_service_RunList(t *testing.T) {
	ctx := t.Context()

	const testAccountUID = uid.AccountUID("test-account-uid")

	tests := []struct {
		name    string
		setup   func(*sdkmocks.MockGlossary)
		params  ListParams
		wantErr bool
		check   func(*testing.T, ListOutput)
	}{
		{
			name:    "validation error — empty account UID",
			setup:   func(m *sdkmocks.MockGlossary) {},
			params:  ListParams{},
			wantErr: true,
		},
		{
			name: "GetByName API error",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "").
					Return(nil, errors.New("API error"))
			},
			params:  ListParams{AccountUID: testAccountUID},
			wantErr: true,
		},
		{
			name: "empty result",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "").
					Return([]glossaryapi.ReadGlossaryResponse{}, nil)
			},
			params: ListParams{AccountUID: testAccountUID},
			check: func(t *testing.T, got ListOutput) {
				if len(got.Glossaries) != 0 {
					t.Errorf("Glossaries len = %v, want 0", len(got.Glossaries))
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name: "returns all glossaries when no name filter",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "").
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: "uid-1", Name: "Alpha", Description: "desc A", LocaleIDs: []string{"en-US"}},
						{GlossaryUid: "uid-2", Name: "Beta", Description: "desc B", LocaleIDs: []string{"en-US", "es-ES"}},
					}, nil)
			},
			params: ListParams{AccountUID: testAccountUID},
			check: func(t *testing.T, got ListOutput) {
				if len(got.Glossaries) != 2 {
					t.Fatalf("Glossaries len = %v, want 2", len(got.Glossaries))
				}
				if got.Glossaries[0].GlossaryUID != "uid-1" {
					t.Errorf("Glossaries[0].GlossaryUID = %v, want uid-1", got.Glossaries[0].GlossaryUID)
				}
				if got.Glossaries[1].Name != "Beta" {
					t.Errorf("Glossaries[1].Name = %v, want Beta", got.Glossaries[1].Name)
				}
				if len(got.Glossaries[1].LocaleIDs) != 2 {
					t.Errorf("Glossaries[1].LocaleIDs len = %v, want 2", len(got.Glossaries[1].LocaleIDs))
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name: "name filter is forwarded to API",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "Marketing").
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: "uid-3", Name: "Marketing"},
					}, nil)
			},
			params: ListParams{AccountUID: testAccountUID, Name: "Marketing"},
			check: func(t *testing.T, got ListOutput) {
				if len(got.Glossaries) != 1 {
					t.Fatalf("Glossaries len = %v, want 1", len(got.Glossaries))
				}
				if got.Glossaries[0].GlossaryUID != "uid-3" {
					t.Errorf("GlossaryUID = %v, want uid-3", got.Glossaries[0].GlossaryUID)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m)
			got, err := service{glossaryApi: m}.RunList(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
