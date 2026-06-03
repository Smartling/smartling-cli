package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Smartling/smartling-cli/services/glossary/glossaryresolver"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/glossary"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

const (
	// TbxExportFileType is the TBX export file type.
	TbxExportFileType = "tbx"
	// defaultExportPageLimit is the page size we send when the caller didn't
	// supply --filter-paging-limit.
	defaultExportPageLimit = 5000
)

var (
	// AllowedExportFileTypes is the set of file types values accepted by the
	// Smartling Glossary Export API.
	AllowedExportFileTypes = []string{"csv", "xlsx", TbxExportFileType}
	// AllowedExportTbxVersions is the set of `tbxVersion` values accepted by the
	// Smartling Glossary Export API when file type is "tbx".
	AllowedExportTbxVersions = []string{"v2", "v3"}
)

func (s service) RunExport(ctx context.Context, params ExportParams) (ExportOutput, error) {
	if err := params.Validate(); err != nil {
		return ExportOutput{}, fmt.Errorf("invalid export params: %w", err)
	}

	glossaryUID, err := glossaryresolver.GetGlossaryUID(ctx, s.glossaryApi, params.AccountUID, params.GlossaryUIDOrName)
	if err != nil {
		return ExportOutput{}, fmt.Errorf("failed to get glossary UID: %w", err)
	}

	// localeIds is required by the export API; when the caller didn't pass any
	// --locale flags, fall back to the glossary's full locale list.
	if len(params.LocaleIDs) == 0 {
		gl, err := s.glossaryApi.Get(ctx, params.AccountUID, glossaryUID)
		if err != nil {
			return ExportOutput{}, fmt.Errorf("get glossary %q to default locale list: %w", glossaryUID, err)
		}
		params.LocaleIDs = gl.LocaleIDs
	}

	resp, err := s.glossaryApi.Export(ctx, params.AccountUID, glossaryUID, toApiExportGlossaryRequest(params))
	if err != nil {
		return ExportOutput{}, fmt.Errorf("failed to get api export glossary: %w", err)
	}
	defer func() {
		if err := resp.Data.Close(); err != nil {
			rlog.Errorf("failed to close export response body: %v", err)
		}
	}()

	outFile := params.OutFile
	if outFile == "" {
		outFile = defaultExportFilename(params.FileType, glossaryUID)
	}
	f, err := os.Create(outFile)
	if err != nil {
		return ExportOutput{}, fmt.Errorf("failed to create export file %q: %w", outFile, err)
	}
	written, copyErr := io.Copy(f, resp.Data)
	if closeErr := f.Close(); closeErr != nil && copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil {
		if rmErr := os.Remove(outFile); rmErr != nil {
			rlog.Errorf("failed to remove partial export file %q: %v", outFile, rmErr)
		}
		return ExportOutput{}, fmt.Errorf("failed to write export file %q: %w", outFile, copyErr)
	}

	return toExportOutput(glossaryUID, outFile, params.FileType, resp, uint64(written)), nil
}

// toAPITbxVersion maps the short CLI/config values ("v2", "v3") to the
// long enum values the Smartling export API expects
// (TBXcoreStructV02 / TBXcoreStructV03). Empty or already-long values pass
// through unchanged so users can supply the canonical form too.
func toAPITbxVersion(v string) string {
	switch strings.ToLower(v) {
	case "v2":
		return "TBXcoreStructV02"
	case "v3":
		return "TBXcoreStructV03"
	default:
		return v
	}
}

// defaultExportFilename builds the fallback output filename when neither
// --out-file nor a server Content-Disposition is available. It uses the
// resolved glossary UID and the requested file type as the extension.
func defaultExportFilename(fileType, glossaryUID string) string {
	ext := strings.ToLower(fileType)
	if ext == "" {
		return glossaryUID
	}
	return glossaryUID + "." + ext
}

// ExportFilter mirrors the Smartling Glossary Export API `filter` object.
// Paging and sorting are intentionally omitted: the /entries/download endpoint
// ignores them (it streams the full filtered result in a fixed order).
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
	LastModifiedBy             LastModifiedBy
	LastModified               LastModified
	CreatedBy                  CreatedBy
	Created                    Created
}

// Created filters entries by creation.
type Created struct {
	Level string
	Type  string
	Date  time.Time
}

// CreatedBy filters entries by creator.
type CreatedBy struct {
	Level   string
	UserIDs []string
}

// LastModified filters entries by last modification.
type LastModified struct {
	Level string
	Type  string
	Date  time.Time
}

// LastModifiedBy filters entries by last modifier.
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
		return fmt.Errorf("file type is required: allowed values are %v", AllowedExportFileTypes)
	}
	fileType := strings.ToLower(p.FileType)
	if !slices.Contains(AllowedExportFileTypes, fileType) {
		return fmt.Errorf("unsupported file type %q: allowed values are %v", p.FileType, AllowedExportFileTypes)
	}
	if fileType == TbxExportFileType {
		if p.TbxVersion == "" {
			return fmt.Errorf("tbx version is required when file type is %q", TbxExportFileType)
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
	BytesWritten uint64
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
		Format:        strings.ToUpper(params.FileType),
		TbxVersion:    toAPITbxVersion(params.TbxVersion),
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
	req.Filter.DntTermSet = f.DntTermSet

	if f.LabelsType != "" {
		req.Filter.Labels = &api.ExportGlossaryLabelsFilter{Type: f.LabelsType}
	}
	if f.Created.Level != "" || f.Created.Type != "" {
		req.Filter.Created = &api.ExportGlossaryDateFilter{
			Level: f.Created.Level,
			Type:  f.Created.Type,
			Date:  f.Created.Date,
		}
	}
	if f.LastModified.Level != "" || f.LastModified.Type != "" {
		req.Filter.LastModified = &api.ExportGlossaryDateFilter{
			Level: f.LastModified.Level,
			Type:  f.LastModified.Type,
			Date:  f.LastModified.Date,
		}
	}
	if f.CreatedBy.Level != "" || len(f.CreatedBy.UserIDs) > 0 {
		req.Filter.CreatedBy = &api.ExportGlossaryUserFilter{
			Level:   f.CreatedBy.Level,
			UserIds: f.CreatedBy.UserIDs,
		}
	}
	if f.LastModifiedBy.Level != "" || len(f.LastModifiedBy.UserIDs) > 0 {
		req.Filter.LastModifiedBy = &api.ExportGlossaryUserFilter{
			Level:   f.LastModifiedBy.Level,
			UserIds: f.LastModifiedBy.UserIDs,
		}
	}
	// limit=0 means "return 0 entries" on Smartling's pagination
	req.Filter.Paging.Limit = defaultExportPageLimit

	return req
}

func toExportOutput(glossaryUID, outFile, fileType string, resp api.ExportGlossaryResponse, bytesWritten uint64) ExportOutput {
	res := ExportOutput{
		GlossaryUID:  glossaryUID,
		OutFile:      outFile,
		FileType:     fileType,
		BytesWritten: bytesWritten,
	}

	summary := struct {
		GlossaryUID  string `json:"glossaryUid"`
		OutFile      string `json:"outFile"`
		FileType     string `json:"fileType"`
		ContentType  string `json:"contentType"`
		BytesWritten uint64 `json:"bytesWritten"`
	}{
		GlossaryUID:  glossaryUID,
		OutFile:      outFile,
		FileType:     fileType,
		ContentType:  resp.ContentType,
		BytesWritten: bytesWritten,
	}
	b, err := json.Marshal(summary)
	if err != nil {
		rlog.Errorf("failed to marshal export output to JSON: %v", err)
		return res
	}
	res.JSON = b
	return res
}
