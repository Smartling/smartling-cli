package files

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	sdk "github.com/Smartling/api-sdk-go"
	sdkjob "github.com/Smartling/api-sdk-go/api/job"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"
)

// recordingAPIClient is a hand-rolled test double for sdk.APIClient. Methods
// not overridden here will nil-panic via the embedded interface, which is
// exactly what we want so unexpected calls show up as test failures.
type recordingAPIClient struct {
	sdk.APIClient
	getStatus           func(fileURI string) (*sdkfile.FileStatus, error)
	downloadTranslation int32
}

func (r *recordingAPIClient) GetFileStatus(_ context.Context, _, fileURI string) (*sdkfile.FileStatus, error) {
	return r.getStatus(fileURI)
}

func (r *recordingAPIClient) DownloadTranslation(_ context.Context, _, _ string, _ sdk.FileDownloadRequest) (io.ReadCloser, error) {
	atomic.AddInt32(&r.downloadTranslation, 1)
	return io.NopCloser(strings.NewReader("translated content")), nil
}

func newServiceWithJobAPI(t *testing.T, projectID string) (service, *jobmocks.MockJob) {
	t.Helper()
	mockJob := jobmocks.NewMockJob(t)
	s := service{
		JobApi: mockJob,
		Config: config.Config{ProjectID: projectID},
	}
	return s, mockJob
}

func TestEnumerateJobFiles_HappyPath(t *testing.T) {
	s, mockJob := newServiceWithJobAPI(t, "proj-1")

	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR", "de-DE"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items: []sdkjob.JobFile{
				{FileURI: "/a.json"},
				{FileURI: "/b.xml"},
			},
			TotalCount: 2,
		}, nil)

	files, locales, err := s.enumerateJobFiles(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("len(files) = %d, want 2", len(files))
	}
	if files[0].FileURI != "/a.json" || files[1].FileURI != "/b.xml" {
		t.Errorf("files = %v, want /a.json then /b.xml", files)
	}
	want := []string{"fr-FR", "de-DE"}
	if !reflect.DeepEqual(locales, want) {
		t.Errorf("locales = %v, want %v", locales, want)
	}
}

func TestEnumerateJobFiles_NotFound(t *testing.T) {
	s, mockJob := newServiceWithJobAPI(t, "proj-1")

	// GetJob and ListFiles run in parallel; for a missing job the underlying
	// API returns 404 for both, which the SDK maps to ErrNotFound.
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "missing").
		Return(sdkjob.GetJobResponse{}, sdkjob.ErrNotFound)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "missing", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{}, sdkjob.ErrNotFound)

	_, _, err := s.enumerateJobFiles(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "missing") || !strings.Contains(msg, "proj-1") {
		t.Errorf("err %q must contain both job UID and project ID", msg)
	}
}

func TestEnumerateJobFiles_EmptyFiles(t *testing.T) {
	s, mockJob := newServiceWithJobAPI(t, "proj-1")

	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{}, nil)

	// enumerateJobFiles itself does not error on empty results — the
	// centralized check in RunPull catches that case. See
	// TestRunPull_JobWithNoFiles_ReturnsError for the end-to-end behavior.
	files, locales, err := s.enumerateJobFiles(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected empty files slice, got %v", files)
	}
	if !reflect.DeepEqual(locales, []string{"fr-FR"}) {
		t.Errorf("locales = %v, want [fr-FR]", locales)
	}
}

func TestRunPull_JobWithNoFiles_ReturnsError(t *testing.T) {
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-empty-files").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-empty-files-uid", JobName: "job-empty-files", TargetLocaleIDs: []string{"fr-FR"}},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-empty-files-uid").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-empty-files-uid",
			TargetLocaleIDs:   []string{"fr-FR"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-empty-files-uid", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{}, nil)

	s := service{
		APIClient: nil,
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1"},
	}

	err := s.RunPull(context.Background(), PullParams{JobUIDOrName: "job-empty-files"})
	if err == nil {
		t.Fatal("expected error for job with no files, got nil")
	}
}

func TestEnumerateJobFiles_ListFilesError(t *testing.T) {
	s, mockJob := newServiceWithJobAPI(t, "proj-1")

	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{TargetLocaleIDs: []string{"fr-FR"}}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{}, errors.New("network blew up"))

	_, _, err := s.enumerateJobFiles(context.Background(), "job-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "network blew up") {
		t.Errorf("err = %v, want wrapped network error", err)
	}
}

func TestRunPull_DryRun_DoesNotCallAPIClient(t *testing.T) {
	// APIClient is intentionally nil — if RunPull ever reaches GetFileStatus
	// or DownloadFile in dry-run mode this test will nil-panic and fail.
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-1").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-1", JobName: "job-1"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR", "de-DE"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items:      []sdkjob.JobFile{{FileURI: "/a.json"}},
			TotalCount: 1,
		}, nil)

	s := service{
		APIClient: nil,
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1"},
	}

	err := s.RunPull(context.Background(), PullParams{
		JobUIDOrName: "job-1",
		DryRun:       true,
	})
	if err != nil {
		t.Fatalf("RunPull dry-run returned error: %v", err)
	}
}

func TestRunPull_EmptyLocaleIntersection_Errors(t *testing.T) {
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-1").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-1", JobName: "job-1"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items:      []sdkjob.JobFile{{FileURI: "/a.json"}},
			TotalCount: 1,
		}, nil)

	s := service{
		JobApi: mockJob,
		Config: config.Config{ProjectID: "proj-1"},
	}

	err := s.RunPull(context.Background(), PullParams{
		JobUIDOrName: "job-1",
		Locales:      []string{"ja-JP"},
		DryRun:       true,
	})
	if err == nil {
		t.Fatal("expected error from empty locale intersection")
	}
	if !strings.Contains(err.Error(), "job-1") {
		t.Errorf("error %q missing job UID context", err)
	}
}

func TestRunPull_JobUIDPlusURI_FiltersJobFiles(t *testing.T) {
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-1").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-1", JobName: "job-1"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items: []sdkjob.JobFile{
				{FileURI: "a.json"},
				{FileURI: "b.xml"},
				{FileURI: "nested/c.json"},
			},
			TotalCount: 3,
		}, nil)

	s := service{
		APIClient: nil, // dry-run: APIClient is unused
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1"},
	}

	// Capture dry-run output to assert exactly which files were selected.
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	err := s.RunPull(context.Background(), PullParams{
		JobUIDOrName: "job-1",
		URI:          "**.json",
		DryRun:       true,
	})

	w.Close()
	os.Stdout = stdout
	out, _ := io.ReadAll(r)

	if err != nil {
		t.Fatalf("RunPull error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 paths (a.json + nested/c.json), got %d: %v", len(lines), lines)
	}
	for _, line := range lines {
		if !strings.HasSuffix(line, ".json") {
			t.Errorf("unexpected non-json line in dry-run output: %q", line)
		}
	}
}

func TestRunPull_JobWithNoLocales_ReturnsError(t *testing.T) {
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-empty").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-empty", JobName: "job-empty"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-empty").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-empty",
			TargetLocaleIDs:   nil,
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-empty", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items:      []sdkjob.JobFile{{FileURI: "/a.json"}},
			TotalCount: 1,
		}, nil)

	s := service{
		APIClient: nil,
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1"},
	}

	err := s.RunPull(context.Background(), PullParams{JobUIDOrName: "job-empty"})
	if err == nil {
		t.Fatal("expected error for job with no target locales, got nil")
	}
	if !strings.Contains(err.Error(), "job-empty") {
		t.Errorf("error %q missing job UID context", err)
	}
}

func TestRunPull_JobUIDPlusURIWithNoMatch_ReturnsError(t *testing.T) {
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-1").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-1", JobName: "job-1"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items: []sdkjob.JobFile{
				{FileURI: "a.json"},
				{FileURI: "b.xml"},
			},
			TotalCount: 2,
		}, nil)

	s := service{
		APIClient: nil,
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1"},
	}

	err := s.RunPull(context.Background(), PullParams{
		JobUIDOrName: "job-1",
		URI:          "**.no_such_extension",
		DryRun:       true,
	})
	if err == nil {
		t.Fatal("expected error when URI glob filters out every job file, got nil")
	}
}

func TestRunPull_Resume_SkipsExistingFiles(t *testing.T) {
	tmpDir := t.TempDir()

	api := &recordingAPIClient{
		getStatus: func(_ string) (*sdkfile.FileStatus, error) {
			return &sdkfile.FileStatus{
				TotalStringCount: 100,
				Items: []sdkfile.FileStatusTranslation{
					{LocaleID: "fr-FR", CompletedStringCount: 100},
					{LocaleID: "de-DE", CompletedStringCount: 100},
				},
			}, nil
		},
	}
	mockJob := jobmocks.NewMockJob(t)
	mockJob.EXPECT().
		SearchByName(context.Background(), "", "job-1").
		Return([]sdkjob.GetJobResponse{
			{TranslationJobUID: "job-1", JobName: "job-1"},
		}, nil)
	mockJob.EXPECT().
		GetJob(context.Background(), "proj-1", "job-1").
		Return(sdkjob.GetJobResponse{
			TranslationJobUID: "job-1",
			TargetLocaleIDs:   []string{"fr-FR", "de-DE"},
		}, nil)
	mockJob.EXPECT().
		ListFiles(context.Background(), "proj-1", "job-1", uint32(500), uint32(0)).
		Return(sdkjob.ListJobFilesResponse{
			Items:      []sdkjob.JobFile{{FileURI: "a.json"}},
			TotalCount: 1,
		}, nil)

	// Pre-create one of the two target paths. The default job format is
	// {{.JobUID}}/{{.Locale}}/{{.FileURI}}, with --directory tmpDir.
	existing := filepath.Join(tmpDir, "job-1", "fr-FR", "a.json")
	if err := os.MkdirAll(filepath.Dir(existing), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.WriteFile(existing, []byte("pre-existing"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	s := service{
		APIClient: api,
		JobApi:    mockJob,
		Config:    config.Config{ProjectID: "proj-1", Threads: 1},
	}

	err := s.RunPull(context.Background(), PullParams{
		JobUIDOrName: "job-1",
		Directory:    tmpDir,
		Resume:       true,
	})
	if err != nil {
		t.Fatalf("RunPull error: %v", err)
	}

	// fr-FR was pre-created → DownloadTranslation should NOT have been called for it.
	// de-DE was not pre-created → DownloadTranslation SHOULD have been called.
	got := atomic.LoadInt32(&api.downloadTranslation)
	if got != 1 {
		t.Errorf("DownloadTranslation call count = %d, want 1 (fr-FR skipped via --resume, de-DE downloaded)", got)
	}

	// Confirm the pre-existing file was untouched.
	body, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("read existing: %v", err)
	}
	if string(body) != "pre-existing" {
		t.Errorf("existing file overwritten; got %q, want %q", body, "pre-existing")
	}
}

func TestFilterLocales(t *testing.T) {
	tests := []struct {
		name        string
		jobLocales  []string
		userLocales []string
		want        []string
	}{
		{
			name:        "user empty returns job locales unchanged",
			jobLocales:  []string{"fr-FR", "de-DE"},
			userLocales: nil,
			want:        []string{"fr-FR", "de-DE"},
		},
		{
			name:        "filter keeps order of jobLocales",
			jobLocales:  []string{"fr-FR", "de-DE", "es-ES"},
			userLocales: []string{"es-ES", "fr-FR"},
			want:        []string{"fr-FR", "es-ES"},
		},
		{
			name:        "case-insensitive match",
			jobLocales:  []string{"fr-FR", "DE-de"},
			userLocales: []string{"FR-fr", "de-DE"},
			want:        []string{"fr-FR", "DE-de"},
		},
		{
			name:        "no overlap returns nil",
			jobLocales:  []string{"fr-FR"},
			userLocales: []string{"ja-JP"},
			want:        nil,
		},
		{
			name:        "user locale not in job is filtered out",
			jobLocales:  []string{"fr-FR", "de-DE"},
			userLocales: []string{"fr-FR", "ja-JP"},
			want:        []string{"fr-FR"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterLocales(tt.jobLocales, tt.userLocales)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterLocales(%v, %v) = %v, want %v",
					tt.jobLocales, tt.userLocales, got, tt.want)
			}
		})
	}
}

func TestPullParams_validate(t *testing.T) {
	tests := []struct {
		name    string
		params  PullParams
		wantErr bool
	}{
		{
			name:    "nothing set is rejected",
			params:  PullParams{},
			wantErr: true,
		},
		{
			name:    "uri alone is accepted",
			params:  PullParams{URI: "foo.json"},
			wantErr: false,
		},
		{
			name:    "all alone is accepted",
			params:  PullParams{All: true},
			wantErr: false,
		},
		{
			name:    "job UID alone is accepted",
			params:  PullParams{JobUIDOrName: "jobUid-1"},
			wantErr: false,
		},
		{
			name:    "uri + all is rejected",
			params:  PullParams{URI: "foo.json", All: true},
			wantErr: true,
		},
		{
			name:    "uri + job-uid is accepted (uri filters job files)",
			params:  PullParams{URI: "**/*.json", JobUIDOrName: "jobUid-1"},
			wantErr: false,
		},
		{
			name:    "all + job-uid is accepted (job-uid wins)",
			params:  PullParams{All: true, JobUIDOrName: "jobUid-1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.params.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPullParams_setDefaultFormatIfEmpty(t *testing.T) {
	tests := []struct {
		name       string
		params     PullParams
		wantFormat string
	}{
		{
			name:       "empty format with job-uid picks job format",
			params:     PullParams{JobUIDOrName: "job-1"},
			wantFormat: format.DefaultFilePullJobFormat,
		},
		{
			name:       "empty format without job-uid picks regular pull format",
			params:     PullParams{URI: "**.json"},
			wantFormat: format.DefaultFilePullFormat,
		},
		{
			name:       "empty format with --all picks regular pull format",
			params:     PullParams{All: true},
			wantFormat: format.DefaultFilePullFormat,
		},
		{
			name:       "user-provided format wins over job-uid default",
			params:     PullParams{JobUIDOrName: "job-1", Format: "{{.FileURI}}.custom"},
			wantFormat: "{{.FileURI}}.custom",
		},
		{
			name:       "user-provided format wins for plain pull",
			params:     PullParams{URI: "**.json", Format: "{{.Locale}}/{{.FileURI}}"},
			wantFormat: "{{.Locale}}/{{.FileURI}}",
		},
		{
			name:       "uri + job-uid still picks job format when no explicit format",
			params:     PullParams{URI: "**.json", JobUIDOrName: "job-1"},
			wantFormat: format.DefaultFilePullJobFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.setDefaultFormatIfEmpty()
			if tt.params.Format != tt.wantFormat {
				t.Errorf("Format = %q, want %q", tt.params.Format, tt.wantFormat)
			}
		})
	}
}
