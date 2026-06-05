package jobsfiles

import (
	"context"
	"errors"
	"testing"

	filesdkmocks "github.com/Smartling/smartling-cli/services/jobs/files/sdkmocks"
	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/file"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// stubList returns a ListFilesFn that yields the given URIs as project files.
func stubList(uris ...string) func(context.Context, string, sdkfile.FilesListRequest) ([]sdkfile.File, error) {
	return func(context.Context, string, sdkfile.FilesListRequest) ([]sdkfile.File, error) {
		files := make([]sdkfile.File, 0, len(uris))
		for _, u := range uris {
			files = append(files, sdkfile.File{FileURI: u})
		}
		return files, nil
	}
}

const (
	addProjectUID = "test-project-id"
	addJobUID     = "aabbccdd1122"
)

func TestRunAdd(t *testing.T) {
	ctx := context.Background()

	t.Run("validation error — no patterns", func(t *testing.T) {
		s := service{jobFile: filesdkmocks.NewMockJobFile(t), job: jobsdkmocks.NewMockJob(t), listFiles: stubList()}
		if _, err := s.RunAdd(ctx, AddParams{ProjectID: addProjectUID, JobUIDOrName: addJobUID}); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("expands pattern and adds each resolved file", func(t *testing.T) {
		job := jobsdkmocks.NewMockJob(t)
		job.EXPECT().GetJob(ctx, addProjectUID, addJobUID).Return(jobapi.GetJobResponse{TranslationJobUID: addJobUID}, nil)
		file := filesdkmocks.NewMockJobFile(t)
		file.EXPECT().Add(ctx, addProjectUID, addJobUID, api.AddRequest{FileURI: "a.json", TargetLocaleIDs: []string{"fr-FR"}}).
			Return(api.Result{SuccessCount: 1}, nil)
		file.EXPECT().Add(ctx, addProjectUID, addJobUID, api.AddRequest{FileURI: "b.json", TargetLocaleIDs: []string{"fr-FR"}}).
			Return(api.Result{SuccessCount: 2, FailCount: 1}, nil)

		s := service{jobFile: file, job: job, listFiles: stubList("a.json", "b.json", "c.txt")}
		got, err := s.RunAdd(ctx, AddParams{
			ProjectID: addProjectUID, JobUIDOrName: addJobUID,
			FilePatterns: []string{"*.json"}, TargetLocaleIDs: []string{"fr-FR"},
		})
		if err != nil {
			t.Fatalf("RunAdd: %v", err)
		}
		if got.SuccessCount != 3 || got.FailCount != 1 {
			t.Errorf("totals = %d/%d, want 3/1", got.SuccessCount, got.FailCount)
		}
		if len(got.Files) != 2 || got.Files[0].FileURI != "a.json" || got.Files[1].FileURI != "b.json" {
			t.Errorf("files = %+v, want a.json,b.json", got.Files)
		}
		if len(got.Unmatched) != 0 {
			t.Errorf("unmatched = %v, want none", got.Unmatched)
		}
	})

	t.Run("unmatched pattern makes no API calls", func(t *testing.T) {
		job := jobsdkmocks.NewMockJob(t)
		job.EXPECT().GetJob(ctx, addProjectUID, addJobUID).Return(jobapi.GetJobResponse{TranslationJobUID: addJobUID}, nil)
		// MockJobFile with no expectations — fails if Add is called.
		s := service{jobFile: filesdkmocks.NewMockJobFile(t), job: job, listFiles: stubList("a.json")}
		got, err := s.RunAdd(ctx, AddParams{ProjectID: addProjectUID, JobUIDOrName: addJobUID, FilePatterns: []string{"*.xml"}})
		if err != nil {
			t.Fatalf("RunAdd: %v", err)
		}
		if len(got.Files) != 0 || len(got.Unmatched) != 1 || got.Unmatched[0] != "*.xml" {
			t.Errorf("got files=%v unmatched=%v, want 0 files / [*.xml]", got.Files, got.Unmatched)
		}
	})

	t.Run("per-file API error is recorded, others continue", func(t *testing.T) {
		job := jobsdkmocks.NewMockJob(t)
		job.EXPECT().GetJob(ctx, addProjectUID, addJobUID).Return(jobapi.GetJobResponse{TranslationJobUID: addJobUID}, nil)
		file := filesdkmocks.NewMockJobFile(t)
		file.EXPECT().Add(ctx, addProjectUID, addJobUID, api.AddRequest{FileURI: "a.json"}).
			Return(api.Result{}, errors.New("boom"))
		file.EXPECT().Add(ctx, addProjectUID, addJobUID, api.AddRequest{FileURI: "b.json"}).
			Return(api.Result{SuccessCount: 1}, nil)

		s := service{jobFile: file, job: job, listFiles: stubList("a.json", "b.json")}
		got, err := s.RunAdd(ctx, AddParams{ProjectID: addProjectUID, JobUIDOrName: addJobUID, FilePatterns: []string{"*.json"}})
		if err != nil {
			t.Fatalf("RunAdd: %v", err)
		}
		if len(got.Files) != 2 || got.Files[0].Error == "" || got.Files[1].Error != "" {
			t.Errorf("files = %+v, want first errored, second ok", got.Files)
		}
		if got.SuccessCount != 1 {
			t.Errorf("SuccessCount = %d, want 1", got.SuccessCount)
		}
	})
}
