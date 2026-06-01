package glossary

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

func TestExportParams_Validate(t *testing.T) {
	validAccount := uid.AccountUID("abc123")

	tests := []struct {
		name    string
		params  ExportParams
		wantErr bool
	}{
		{
			name: "valid csv",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "csv",
			},
			wantErr: false,
		},
		{
			name: "valid xlsx",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "xlsx",
			},
			wantErr: false,
		},
		{
			name: "valid tbx with version",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "tbx",
				TbxVersion:        "v2",
			},
			wantErr: false,
		},
		{
			name: "case insensitive file type",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "CSV",
			},
			wantErr: false,
		},
		{
			name:    "empty account UID",
			params:  ExportParams{GlossaryUIDOrName: "my-glossary", FileType: "csv"},
			wantErr: true,
		},
		{
			name:    "empty glossary UID or name",
			params:  ExportParams{AccountUID: validAccount, FileType: "csv"},
			wantErr: true,
		},
		{
			name:    "empty file type",
			params:  ExportParams{AccountUID: validAccount, GlossaryUIDOrName: "my-glossary"},
			wantErr: true,
		},
		{
			name: "unsupported file type",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "docx",
			},
			wantErr: true,
		},
		{
			name: "tbx without version",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "tbx",
			},
			wantErr: true,
		},
		{
			name: "tbx with invalid version",
			params: ExportParams{
				AccountUID:        validAccount,
				GlossaryUIDOrName: "my-glossary",
				FileType:          "tbx",
				TbxVersion:        "v9",
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

func Test_service_RunExport(t *testing.T) {
	ctx := t.Context()

	const (
		testAccountUID   = uid.AccountUID("test-account-uid")
		testGlossaryUID  = "00000000-0000-0000-0000-000000000001"
		testGlossaryName = "my-glossary"
	)

	setupGetByName := func(m *sdkmocks.MockGlossary) {
		m.EXPECT().GetByName(ctx, testAccountUID, testGlossaryName).
			Return([]glossaryapi.ReadGlossaryResponse{
				{GlossaryUid: testGlossaryUID, Name: testGlossaryName},
			}, nil)
	}

	baseParams := func(outFile string) ExportParams {
		return ExportParams{
			AccountUID:        testAccountUID,
			GlossaryUIDOrName: testGlossaryName,
			FileType:          "csv",
			LocaleIDs:         []string{"en-US"},
			OutFile:           outFile,
		}
	}

	csvResponse := func() glossaryapi.ExportGlossaryResponse {
		return glossaryapi.ExportGlossaryResponse{
			ContentType: "text/csv",
			Data:        io.NopCloser(strings.NewReader("term,translation\nhello,hola")),
		}
	}

	tests := []struct {
		name    string
		setup   func(*sdkmocks.MockGlossary, string)
		params  func(outFile string) ExportParams
		wantErr bool
		check   func(*testing.T, ExportOutput)
	}{
		{
			name: "GetGlossaryUID error — glossary not found",
			setup: func(m *sdkmocks.MockGlossary, _ string) {
				m.EXPECT().GetByName(ctx, testAccountUID, testGlossaryName).
					Return(nil, errors.New("network error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name:  "validation error — empty FileType",
			setup: func(m *sdkmocks.MockGlossary, _ string) {},
			params: func(outFile string) ExportParams {
				p := baseParams(outFile)
				p.FileType = ""
				return p
			},
			wantErr: true,
		},
		{
			name: "Export API error",
			setup: func(m *sdkmocks.MockGlossary, outFile string) {
				setupGetByName(m)
				p := baseParams(outFile)
				m.EXPECT().Export(ctx, testAccountUID, testGlossaryUID, toApiExportGlossaryRequest(p)).
					Return(glossaryapi.ExportGlossaryResponse{}, errors.New("export API error"))
			},
			params:  baseParams,
			wantErr: true,
		},
		{
			name: "success with explicit locale IDs",
			setup: func(m *sdkmocks.MockGlossary, outFile string) {
				setupGetByName(m)
				p := baseParams(outFile)
				m.EXPECT().Export(ctx, testAccountUID, testGlossaryUID, toApiExportGlossaryRequest(p)).
					Return(csvResponse(), nil)
			},
			params: baseParams,
			check: func(t *testing.T, got ExportOutput) {
				if got.GlossaryUID != testGlossaryUID {
					t.Errorf("GlossaryUID = %v, want %v", got.GlossaryUID, testGlossaryUID)
				}
				if got.FileType != "csv" {
					t.Errorf("FileType = %v, want csv", got.FileType)
				}
				if got.BytesWritten == 0 {
					t.Error("BytesWritten should be > 0")
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name: "success — no locale IDs defaults to glossary locales",
			setup: func(m *sdkmocks.MockGlossary, outFile string) {
				setupGetByName(m)
				m.EXPECT().Get(ctx, testAccountUID, testGlossaryUID).
					Return(glossaryapi.ReadGlossaryResponse{
						GlossaryUid: testGlossaryUID,
						LocaleIDs:   []string{"en-US", "es-ES"},
					}, nil)
				p := ExportParams{
					AccountUID:        testAccountUID,
					GlossaryUIDOrName: testGlossaryName,
					FileType:          "csv",
					LocaleIDs:         []string{"en-US", "es-ES"},
					OutFile:           outFile,
				}
				m.EXPECT().Export(ctx, testAccountUID, testGlossaryUID, toApiExportGlossaryRequest(p)).
					Return(csvResponse(), nil)
			},
			params: func(outFile string) ExportParams {
				p := baseParams(outFile)
				p.LocaleIDs = nil
				return p
			},
			check: func(t *testing.T, got ExportOutput) {
				if got.GlossaryUID != testGlossaryUID {
					t.Errorf("GlossaryUID = %v, want %v", got.GlossaryUID, testGlossaryUID)
				}
				if got.BytesWritten == 0 {
					t.Error("BytesWritten should be > 0")
				}
			},
		},
		{
			name: "Get error when defaulting locale IDs",
			setup: func(m *sdkmocks.MockGlossary, _ string) {
				setupGetByName(m)
				m.EXPECT().Get(ctx, testAccountUID, testGlossaryUID).
					Return(glossaryapi.ReadGlossaryResponse{}, errors.New("get error"))
			},
			params: func(outFile string) ExportParams {
				p := baseParams(outFile)
				p.LocaleIDs = nil
				return p
			},
			wantErr: true,
		},
		{
			name: "default output filename when OutFile is empty",
			setup: func(m *sdkmocks.MockGlossary, _ string) {
				setupGetByName(m)
				p := ExportParams{
					AccountUID:        testAccountUID,
					GlossaryUIDOrName: testGlossaryName,
					FileType:          "csv",
					LocaleIDs:         []string{"en-US"},
				}
				m.EXPECT().Export(ctx, testAccountUID, testGlossaryUID, toApiExportGlossaryRequest(p)).
					Return(csvResponse(), nil)
			},
			params: func(_ string) ExportParams {
				return ExportParams{
					AccountUID:        testAccountUID,
					GlossaryUIDOrName: testGlossaryName,
					FileType:          "csv",
					LocaleIDs:         []string{"en-US"},
				}
			},
			check: func(t *testing.T, got ExportOutput) {
				expected := testGlossaryUID + ".csv"
				if got.OutFile != expected {
					t.Errorf("OutFile = %v, want %v", got.OutFile, expected)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outFile := filepath.Join(t.TempDir(), "export.csv")
			if tt.name == "default output filename when OutFile is empty" {
				t.Chdir(t.TempDir())
				outFile = ""
			}
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m, outFile)
			got, err := service{glossaryApi: m}.RunExport(ctx, tt.params(outFile))
			if (err != nil) != tt.wantErr {
				t.Errorf("RunExport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
