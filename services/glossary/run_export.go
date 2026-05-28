package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Smartling/smartling-cli/services/glossary/glossaryresolver"

	api "github.com/Smartling/api-sdk-go/api/glossary"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

const TbxExportFileType = "tbx"

var (
	// AllowedExportFileTypes is the set of file types values accepted by the
	// Smartling Glossary Export API.
	AllowedExportFileTypes = []string{"csv", "xlsx", TbxExportFileType}
	// AllowedExportTbxVersions is the set of `tbxVersion` values accepted by the
	// Smartling Glossary Export API when file type is "tbx".
	AllowedExportTbxVersions = []string{"v2", "v3"}
)

func (s service) RunExport(ctx context.Context, params ExportParams) (ExportOutput, error) {
	glossaryUID, err := glossaryresolver.GetGlossaryUID(ctx, s.glossaryApi, params.AccountUID, params.GlossaryUIDOrName)
	if err != nil {
		return ExportOutput{}, fmt.Errorf("failed to get glossary UID: %w", err)
	}

	if err := params.Validate(); err != nil {
		return ExportOutput{}, fmt.Errorf("invalid export params: %w", err)
	}

	resp, err := s.glossaryApi.Export(ctx, params.AccountUID, glossaryUID, toApiExportGlossaryRequest(params))
	if err != nil {
		return ExportOutput{}, fmt.Errorf("failed to get api export glossary: %w", err)
	}

	outFile := params.OutFile
	if outFile == "" {
		outFile = resp.Filename
	}
	if err := os.WriteFile(outFile, resp.Data, 0o644); err != nil {
		return ExportOutput{}, fmt.Errorf("failed to write export file %q: %w", outFile, err)
	}

	return toExportOutput(glossaryUID, outFile, params.FileType, resp), nil
}

// ExportFilter mirrors the Smartling Glossary Export API `filter` object.
type ExportFilter struct {
	Query                      string
	LocaleID                   []string
	EntryUIDs                  []string
	EntryState                 string
	MissingTranslationLocaleID string
	PresentTranslationLocaleID string
	DntLocaleID                string
	ReturnFallbackTranslations bool
	LabelsType                 string
	DntTermSet                 bool
	Paging                     Paging
	LastModifiedBy             LastModifiedBy
	LastModified               LastModified
	CreatedBy                  CreatedBy
	Created                    Created
	Sorting                    Sorting
}

type Paging struct {
	Offset int
	Limit  int
}

type Sorting struct {
	Field     string
	Direction string
	LocaleID  string
}

type Created struct {
	Level string
	Type  string
	Date  time.Time
}

type CreatedBy struct {
	Level   string
	UserIDs []string
}

type LastModified struct {
	Level string
	Type  string
	Date  time.Time
}
type LastModifiedBy struct {
	Level   string
	UserIDs []string
}

// ExportParams carries the full glossary-export request from CLI to service.
type ExportParams struct {
	GlossaryUIDOrName string
	AccountUID        uid.AccountUID
	OutFile           string
	FileType          string
	TbxVersion        string
	FocusLocaleID     string
	LocaleIDs         []string
	SkipEntries       bool
	Filter            ExportFilter
}

// Validate checks that ExportParams carry the fields required by the
// Smartling Glossary Export API
// (https://api-reference.smartling.com/#tag/Glossary-API/operation/exportGlossary).
func (p ExportParams) Validate() error {
	if err := p.AccountUID.Validate(); err != nil {
		return err
	}
	if p.GlossaryUIDOrName == "" {
		return smerror.ErrEmptyParam("GlossaryUIDOrName")
	}
	if p.FileType == "" {
		return smerror.ErrEmptyParam("FileType")
	}
	fileType := strings.ToLower(p.FileType)
	if !slices.Contains(AllowedExportFileTypes, fileType) {
		return fmt.Errorf("unsupported file type %q: allowed values are %v", p.FileType, AllowedExportFileTypes)
	}
	if fileType == TbxExportFileType {
		if p.TbxVersion == "" {
			return smerror.ErrEmptyParam("TbxVersion")
		}
		if !slices.Contains(AllowedExportTbxVersions, strings.ToLower(p.TbxVersion)) {
			return fmt.Errorf("unsupported tbx version %q: allowed values are %v", p.TbxVersion, AllowedExportTbxVersions)
		}
	}
	return nil
}

// ExportOutput represents the result of a glossary export.
type ExportOutput struct {
	GlossaryUID  string
	OutFile      string
	FileType     string
	BytesWritten int
	JSON         []byte
}

// JSONBytes returns the raw JSON payload of the export response.
func (e ExportOutput) JSONBytes() []byte { return e.JSON }

// SimpleLines returns a human-readable summary of the export.
func (e ExportOutput) SimpleLines() []string {
	return []string{
		fmt.Sprintf("Glossary UID: %s", e.GlossaryUID),
		fmt.Sprintf("Output file:  %s", e.OutFile),
		fmt.Sprintf("File type:    %s", e.FileType),
		fmt.Sprintf("Bytes written: %d", e.BytesWritten),
	}
}

// TableData returns the export summary with one column per field and a
// single row of values.
func (e ExportOutput) TableData() ([]string, [][]string) {
	headers := []string{"GLOSSARY UID", "OUTPUT FILE", "FILE TYPE", "BYTES WRITTEN"}
	rows := [][]string{
		{e.GlossaryUID, e.OutFile, e.FileType, fmt.Sprintf("%d", e.BytesWritten)},
	}
	return headers, rows
}

func toApiExportGlossaryRequest(params ExportParams) api.ExportGlossaryRequest {
	req := api.ExportGlossaryRequest{
		Format:        params.FileType,
		TbxVersion:    params.TbxVersion,
		FocusLocaleId: params.FocusLocaleID,
		LocaleIds:     params.LocaleIDs,
		SkipEntries:   params.SkipEntries,
	}

	f := params.Filter
	req.Filter.Query = f.Query
	req.Filter.LocaleIds = f.LocaleID
	req.Filter.EntryUids = f.EntryUIDs
	req.Filter.EntryState = f.EntryState
	req.Filter.MissingTranslationLocaleId = f.MissingTranslationLocaleID
	req.Filter.PresentTranslationLocaleId = f.PresentTranslationLocaleID
	req.Filter.DntLocaleId = f.DntLocaleID
	req.Filter.ReturnFallbackTranslations = f.ReturnFallbackTranslations
	req.Filter.Labels.Type = f.LabelsType
	req.Filter.DntTermSet = f.DntTermSet

	req.Filter.Created.Level = f.Created.Level
	req.Filter.Created.Type = f.Created.Type
	req.Filter.Created.Date = f.Created.Date

	req.Filter.LastModified.Level = f.LastModified.Level
	req.Filter.LastModified.Type = f.LastModified.Type
	req.Filter.LastModified.Date = f.LastModified.Date

	req.Filter.CreatedBy.Level = f.CreatedBy.Level
	req.Filter.CreatedBy.UserIds = f.CreatedBy.UserIDs

	req.Filter.LastModifiedBy.Level = f.LastModifiedBy.Level
	req.Filter.LastModifiedBy.UserIds = f.LastModifiedBy.UserIDs

	req.Filter.Paging.Offset = f.Paging.Offset
	req.Filter.Paging.Limit = f.Paging.Limit

	req.Filter.Sorting.Field = f.Sorting.Field
	req.Filter.Sorting.Direction = f.Sorting.Direction
	req.Filter.Sorting.LocaleId = f.Sorting.LocaleID

	return req
}

func toExportOutput(glossaryUID, outFile, fileType string, resp api.ExportGlossaryResponse) ExportOutput {
	out := ExportOutput{
		GlossaryUID:  glossaryUID,
		OutFile:      outFile,
		FileType:     fileType,
		BytesWritten: len(resp.Data),
	}

	summary := struct {
		GlossaryUID  string `json:"glossaryUid"`
		OutFile      string `json:"outFile"`
		FileType     string `json:"fileType"`
		ContentType  string `json:"contentType"`
		BytesWritten int    `json:"bytesWritten"`
	}{
		GlossaryUID:  glossaryUID,
		OutFile:      outFile,
		FileType:     fileType,
		ContentType:  resp.ContentType,
		BytesWritten: len(resp.Data),
	}
	if b, err := json.Marshal(summary); err == nil {
		out.JSON = b
	}

	return out
}
