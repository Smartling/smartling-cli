package glossaryresolver

import (
	"errors"
	"testing"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

const (
	testAccountUID  = uid.AccountUID("test-account-uid")
	testGlossaryUID = "00000000-0000-0000-0000-000000000001"
)

func TestGetGlossaryUID(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name              string
		setup             func(*sdkmocks.MockGlossary)
		glossaryUIDOrName string
		want              string
		wantErr           bool
	}{
		{
			name:              "empty input returns not found",
			setup:             func(m *sdkmocks.MockGlossary) {},
			glossaryUIDOrName: "",
			wantErr:           true,
		},
		{
			name:              "whitespace input returns not found",
			setup:             func(m *sdkmocks.MockGlossary) {},
			glossaryUIDOrName: "   ",
			wantErr:           true,
		},
		{
			name: "UUID input found via Get",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().Get(ctx, string(testAccountUID), testGlossaryUID).
					Return(glossaryapi.ReadGlossaryResponse{GlossaryUid: testGlossaryUID, Name: "My Glossary"}, nil)
			},
			glossaryUIDOrName: testGlossaryUID,
			want:              testGlossaryUID,
		},
		{
			name: "UUID input not found via Get falls through to GetByName",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().Get(ctx, string(testAccountUID), testGlossaryUID).
					Return(glossaryapi.ReadGlossaryResponse{}, glossaryapi.ErrGlossaryNotFound)
				m.EXPECT().GetByName(ctx, string(testAccountUID), testGlossaryUID).
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: testGlossaryUID, Name: testGlossaryUID},
					}, nil)
			},
			glossaryUIDOrName: testGlossaryUID,
			want:              testGlossaryUID,
		},
		{
			name: "UUID input Get returns unexpected error",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().Get(ctx, string(testAccountUID), testGlossaryUID).
					Return(glossaryapi.ReadGlossaryResponse{}, errors.New("network error"))
			},
			glossaryUIDOrName: testGlossaryUID,
			wantErr:           true,
		},
		{
			name: "name input found via GetByName exact match",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "My Glossary").
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: testGlossaryUID, Name: "My Glossary"},
						{GlossaryUid: "aaaaaaaa-0000-0000-0000-000000000000", Name: "Other"},
					}, nil)
			},
			glossaryUIDOrName: "My Glossary",
			want:              testGlossaryUID,
		},
		{
			name: "name input falls back to first result when no exact match",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "partial").
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: testGlossaryUID, Name: "partial match"},
					}, nil)
			},
			glossaryUIDOrName: "partial",
			want:              testGlossaryUID,
		},
		{
			name: "name input GetByName returns empty list",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "unknown").
					Return([]glossaryapi.ReadGlossaryResponse{}, nil)
			},
			glossaryUIDOrName: "unknown",
			wantErr:           true,
		},
		{
			name: "name input first result has empty UID",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "bad").
					Return([]glossaryapi.ReadGlossaryResponse{
						{GlossaryUid: "", Name: "bad"},
					}, nil)
			},
			glossaryUIDOrName: "bad",
			wantErr:           true,
		},
		{
			name: "name input GetByName error",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, string(testAccountUID), "My Glossary").
					Return(nil, errors.New("network error"))
			},
			glossaryUIDOrName: "My Glossary",
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m)
			got, err := GetGlossaryUID(ctx, m, testAccountUID, tt.glossaryUIDOrName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGlossaryUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGlossaryUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
