package glossary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Smartling/smartling-cli/services/glossary/glossaryresolver"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/Smartling/api-sdk-go/api/glossary"
	api "github.com/Smartling/api-sdk-go/api/glossary"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

var (
	errFailedImport             = errors.New("glossary import failed")
	errImportConfirmationFailed = errors.New("glossary import confirmation failed")
	pollingInterval             = time.Second
)

type ImportParams struct {
	GlossaryUIDOrName string
	AccountUID        uid.AccountUID
	ArchiveMode       bool
	ImportFile        ImportFile
}

// Validate enforces the fields required by the Smartling Glossary Import API.
func (p ImportParams) Validate() error {
	if err := p.AccountUID.Validate(); err != nil {
		return err
	}
	if p.GlossaryUIDOrName == "" {
		return smerror.ErrEmptyParam("GlossaryUIDOrName")
	}
	if p.ImportFile.Path == "" {
		return smerror.ErrEmptyParam("ImportFile.Path")
	}
	if p.ImportFile.Name == "" {
		return smerror.ErrEmptyParam("ImportFile.Name")
	}
	if p.ImportFile.MediaType == "" {
		return smerror.ErrEmptyParam("ImportFile.MediaType")
	}
	return nil
}

type ImportFile struct {
	Path string
	Name string
	// MediaType: "text/csv" | "text/xml" (TBX) | "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" — server derives format from this.
	MediaType string
}

// ImportOutput represents the result of a glossary import.
type ImportOutput struct {
	GlossaryUID  string
	ImportUID    string
	ImportStatus string
	SourceFile   string
	EntryChanges api.ImportEntryChanges
	Warnings     []api.ImportWarning
	JSON         []byte
}

// JSONBytes returns the raw JSON payload of the import response.
func (p ImportOutput) JSONBytes() []byte { return p.JSON }

// SimpleLines returns a human-readable summary of the import.
func (p ImportOutput) SimpleLines() []string {
	lines := []string{
		fmt.Sprintf("Glossary UID:  %s", p.GlossaryUID),
		fmt.Sprintf("Import UID:    %s", p.ImportUID),
		fmt.Sprintf("Import status: %s", p.ImportStatus),
		fmt.Sprintf("Source file:   %s", p.SourceFile),
		fmt.Sprintf("New entries:           %d", p.EntryChanges.NewEntries),
		fmt.Sprintf("Existing entry updates: %d", p.EntryChanges.ExistingEntryUpdates),
		fmt.Sprintf("Not matched entries:   %d", p.EntryChanges.NotMatchedEntries),
		fmt.Sprintf("Entries to archive:    %d", p.EntryChanges.EntriesToArchive),
	}
	for _, w := range p.Warnings {
		lines = append(lines, fmt.Sprintf("Warning [%s]: %s", w.Key, w.Message))
	}
	return lines
}

// TableData returns the import summary as a single-row table; warnings are in simple/JSON only.
func (p ImportOutput) TableData() ([]string, [][]string) {
	headers := []string{
		"GLOSSARY UID", "IMPORT UID", "STATUS", "SOURCE FILE",
		"NEW", "UPDATED", "NOT MATCHED", "TO ARCHIVE",
	}
	rows := [][]string{
		{
			p.GlossaryUID,
			p.ImportUID,
			p.ImportStatus,
			p.SourceFile,
			fmt.Sprintf("%d", p.EntryChanges.NewEntries),
			fmt.Sprintf("%d", p.EntryChanges.ExistingEntryUpdates),
			fmt.Sprintf("%d", p.EntryChanges.NotMatchedEntries),
			fmt.Sprintf("%d", p.EntryChanges.EntriesToArchive),
		},
	}
	return headers, rows
}

func (s service) RunImport(ctx context.Context, params ImportParams) (ImportOutput, error) {
	if err := params.Validate(); err != nil {
		return ImportOutput{}, fmt.Errorf("invalid import params: %w", err)
	}
	glossaryUID, err := glossaryresolver.GetGlossaryUID(ctx, s.glossaryApi, params.AccountUID, params.GlossaryUIDOrName)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("failed to get glossary UID: %w", err)
	}

	apiImportGlossaryRequest, err := toApiImportGlossaryRequest(params)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("failed to build import glossary request: %w", err)
	}
	importGlossaryResponse, err := s.glossaryApi.Import(ctx, params.AccountUID, glossaryUID, apiImportGlossaryRequest)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("failed to run glossary import: %w", err)
	}
	importConfirmed, err := s.glossaryApi.ImportConfirm(ctx, params.AccountUID, glossaryUID, importGlossaryResponse.ImportUID)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("failed to confirm glossary import: %w", err)
	}
	if !importConfirmed {
		return ImportOutput{}, errImportConfirmationFailed
	}
	finalResponse := importGlossaryResponse
	for {
		importStatusResponse, err := s.glossaryApi.ImportStatus(ctx, params.AccountUID, glossaryUID, importGlossaryResponse.ImportUID)
		if err != nil {
			return ImportOutput{}, fmt.Errorf("failed to get glossary import status: %w", err)
		}
		if importStatusResponse.ImportStatus == glossary.FailedImportStatus {
			return ImportOutput{}, errFailedImport
		}
		if importStatusResponse.ImportStatus == glossary.SuccessfulImportStatus {
			finalResponse.ImportStatus = importStatusResponse.ImportStatus
			break
		}
		select {
		case <-ctx.Done():
			return ImportOutput{}, ctx.Err()
		case <-time.After(pollingInterval):
		}
	}
	return toImportOutput(glossaryUID, params.ImportFile.Path, finalResponse), nil
}

func toApiImportGlossaryRequest(params ImportParams) (api.ImportGlossaryRequest, error) {
	data, err := os.ReadFile(params.ImportFile.Path)
	if err != nil {
		return api.ImportGlossaryRequest{}, fmt.Errorf("read import file %q: %w", params.ImportFile.Path, err)
	}
	return api.ImportGlossaryRequest{
		File:        data,
		FileName:    params.ImportFile.Name,
		MediaType:   params.ImportFile.MediaType,
		ArchiveMode: params.ArchiveMode,
	}, nil
}

func toImportOutput(glossaryUID, sourceFile string, resp api.ImportGlossaryResponse) ImportOutput {
	res := ImportOutput{
		GlossaryUID:  glossaryUID,
		ImportUID:    resp.ImportUID,
		ImportStatus: resp.ImportStatus,
		SourceFile:   sourceFile,
		EntryChanges: resp.EntryChanges,
		Warnings:     resp.Warnings,
	}

	summary := struct {
		GlossaryUID  string                 `json:"glossaryUid"`
		ImportUID    string                 `json:"importUid"`
		ImportStatus string                 `json:"importStatus"`
		SourceFile   string                 `json:"sourceFile"`
		EntryChanges api.ImportEntryChanges `json:"entryChanges"`
		Warnings     []api.ImportWarning    `json:"warnings,omitempty"`
	}{
		GlossaryUID:  glossaryUID,
		ImportUID:    resp.ImportUID,
		ImportStatus: resp.ImportStatus,
		SourceFile:   sourceFile,
		EntryChanges: resp.EntryChanges,
		Warnings:     resp.Warnings,
	}
	b, err := json.Marshal(summary)
	if err != nil {
		rlog.Errorf("failed to marshal import output to JSON: %v", err)
		return res
	}
	res.JSON = b
	return res
}
