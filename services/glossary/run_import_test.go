package glossary

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

func TestImportParams_Validate(t *testing.T) {
	validAccount := uid.AccountUID("abc123")
	validFile := ImportFile{
		Path:      "/tmp/terms.csv",
		Name:      "terms.csv",
		MediaType: "text/csv",
	}

	tests := []struct {
		name    string
		params  ImportParams
		wantErr bool
	}{
		{
			name: "valid",
			params: ImportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				ImportFile:        validFile,
			},
			wantErr: false,
		},
		{
			name:    "empty account UID",
			params:  ImportParams{GlossaryUIDOrName: "my-glossary", ImportFile: validFile},
			wantErr: true,
		},
		{
			name:    "empty glossary UID or name",
			params:  ImportParams{AccountUID: validAccount, ImportFile: validFile},
			wantErr: true,
		},
		{
			name: "empty file path",
			params: ImportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				ImportFile:        ImportFile{Name: "terms.csv", MediaType: "text/csv"},
			},
			wantErr: true,
		},
		{
			name: "empty file name",
			params: ImportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				ImportFile:        ImportFile{Path: "/tmp/terms.csv", MediaType: "text/csv"},
			},
			wantErr: true,
		},
		{
			name: "empty media type",
			params: ImportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				ImportFile:        ImportFile{Path: "/tmp/terms.csv", Name: "terms.csv"},
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

func Test_service_RunImport(t *testing.T) {
	ctx := t.Context()

	t.Cleanup(func() { pollingInterval = time.Second })
	pollingInterval = 0

	const (
		testAccountUID   = uid.AccountUID("test-account-uid")
		testGlossaryUID  = "00000000-0000-0000-0000-000000000001"
		testGlossaryName = "my-glossary"
		testImportUID    = "import-uid-001"
	)

	fileContent := []byte("term,translation\nhello,hola")

	makeImportFile := func(t *testing.T) ImportFile {
		t.Helper()
		path := filepath.Join(t.TempDir(), "terms.csv")
		if err := os.WriteFile(path, fileContent, 0o600); err != nil {
			t.Fatalf("write temp file: %v", err)
		}
		return ImportFile{Path: path, Name: "terms.csv", MediaType: "text/csv"}
	}

	setupGetByName := func(m *sdkmocks.MockGlossary) {
		m.EXPECT().GetByName(ctx, testAccountUID, testGlossaryName).
			Return([]glossaryapi.GetGlossaryResponse{
				{GlossaryUID: testGlossaryUID, Name: testGlossaryName},
			}, nil)
	}

	baseParams := func(f ImportFile) ImportParams {
		return ImportParams{
			AccountUID:        testAccountUID,
			GlossaryUIDOrName: testGlossaryName,
			ImportFile:        f,
		}
	}

	importResponse := glossaryapi.ImportGlossaryResponse{
		GlossaryUID:  testGlossaryUID,
		ImportUID:    testImportUID,
		ImportStatus: glossaryapi.PendingImportStatus,
		EntryChanges: glossaryapi.ImportEntryChanges{NewEntries: 3, ExistingEntryUpdates: 1},
	}

	expectedReq := func(f ImportFile) glossaryapi.ImportGlossaryRequest {
		return glossaryapi.ImportGlossaryRequest{
			File:      fileContent,
			FileName:  f.Name,
			MediaType: f.MediaType,
		}
	}

	tests := []struct {
		name    string
		setup   func(*sdkmocks.MockGlossary, ImportFile)
		params  func(ImportFile) ImportParams
		wantErr bool
		check   func(*testing.T, ImportOutput)
	}{
		{
			name:  "validation error — empty file path",
			setup: func(m *sdkmocks.MockGlossary, _ ImportFile) {},
			params: func(_ ImportFile) ImportParams {
				return ImportParams{
					AccountUID:        testAccountUID,
					GlossaryUIDOrName: testGlossaryName,
					ImportFile:        ImportFile{Name: "terms.csv", MediaType: "text/csv"},
				}
			},
			wantErr: true,
		},
		{
			name: "GetGlossaryUID error",
			setup: func(m *sdkmocks.MockGlossary, _ ImportFile) {
				m.EXPECT().GetByName(ctx, testAccountUID, testGlossaryName).
					Return(nil, errors.New("network error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "file read error — path does not exist",
			setup: func(m *sdkmocks.MockGlossary, _ ImportFile) {
				setupGetByName(m)
			},
			params: func(_ ImportFile) ImportParams {
				return ImportParams{
					AccountUID:        testAccountUID,
					GlossaryUIDOrName: testGlossaryName,
					ImportFile:        ImportFile{Path: "/nonexistent/terms.csv", Name: "terms.csv", MediaType: "text/csv"},
				}
			},
			wantErr: true,
		},
		{
			name: "Import API error",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(glossaryapi.ImportGlossaryResponse{}, errors.New("import API error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "ImportConfirm error",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(false, errors.New("confirm error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "ImportConfirm returns false",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(false, nil)
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "ImportStatus error",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(true, nil)
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{}, errors.New("status error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "ImportStatus returns FAILED",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(true, nil)
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{ImportStatus: glossaryapi.FailedImportStatus}, nil)
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "success — SUCCESSFUL on first poll",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(true, nil)
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{ImportStatus: glossaryapi.SuccessfulImportStatus}, nil)
			},
			params: baseParams,
			check: func(t *testing.T, got ImportOutput) {
				if got.GlossaryUID != testGlossaryUID {
					t.Errorf("GlossaryUID = %v, want %v", got.GlossaryUID, testGlossaryUID)
				}
				if got.ImportUID != testImportUID {
					t.Errorf("ImportUID = %v, want %v", got.ImportUID, testImportUID)
				}
				if got.ImportStatus != glossaryapi.SuccessfulImportStatus {
					t.Errorf("ImportStatus = %v, want %v", got.ImportStatus, glossaryapi.SuccessfulImportStatus)
				}
				if got.EntryChanges.NewEntries != 3 {
					t.Errorf("NewEntries = %v, want 3", got.EntryChanges.NewEntries)
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name: "success — PENDING then SUCCESSFUL",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, expectedReq(f)).
					Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(true, nil)
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{ImportStatus: glossaryapi.PendingImportStatus}, nil).
					Call.Once()
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{ImportStatus: glossaryapi.SuccessfulImportStatus}, nil).
					Call.Once()
			},
			params: baseParams,
			check: func(t *testing.T, got ImportOutput) {
				if got.ImportStatus != glossaryapi.SuccessfulImportStatus {
					t.Errorf("ImportStatus = %v, want %v", got.ImportStatus, glossaryapi.SuccessfulImportStatus)
				}
			},
		},
		{
			name: "success — archive mode",
			setup: func(m *sdkmocks.MockGlossary, f ImportFile) {
				setupGetByName(m)
				m.EXPECT().Import(ctx, testAccountUID, testGlossaryUID, glossaryapi.ImportGlossaryRequest{
					File:        fileContent,
					FileName:    f.Name,
					MediaType:   f.MediaType,
					ArchiveMode: true,
				}).Return(importResponse, nil)
				m.EXPECT().ImportConfirm(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(true, nil)
				m.EXPECT().ImportStatus(ctx, testAccountUID, testGlossaryUID, testImportUID).
					Return(glossaryapi.ImportStatusResponse{ImportStatus: glossaryapi.SuccessfulImportStatus}, nil)
			},
			params: func(f ImportFile) ImportParams {
				p := baseParams(f)
				p.ArchiveMode = true
				return p
			},
			check: func(t *testing.T, got ImportOutput) {
				if got.ImportStatus != glossaryapi.SuccessfulImportStatus {
					t.Errorf("ImportStatus = %v, want %v", got.ImportStatus, glossaryapi.SuccessfulImportStatus)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := makeImportFile(t)
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m, f)
			got, err := service{glossaryApi: m}.RunImport(ctx, tt.params(f))
			if (err != nil) != tt.wantErr {
				t.Errorf("RunImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
