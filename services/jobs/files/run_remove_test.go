package jobsfiles

import (
	"context"
	"testing"

	filesdkmocks "github.com/Smartling/smartling-cli/services/jobs/files/sdkmocks"
	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/file"
)

func TestRunRemove(t *testing.T) {
	ctx := context.Background()

	t.Run("validation error — no patterns", func(t *testing.T) {
		s := service{jobFile: filesdkmocks.NewMockJobFile(t), job: jobsdkmocks.NewMockJob(t), listFiles: stubList()}
		if _, err := s.RunRemove(ctx, RemoveParams{ProjectID: addProjectUID, JobUIDOrName: addJobUID}); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("expands pattern and removes each resolved file", func(t *testing.T) {
		job := jobsdkmocks.NewMockJob(t)
		job.EXPECT().GetJob(ctx, addProjectUID, addJobUID).Return(jobapi.GetJobResponse{TranslationJobUID: addJobUID}, nil)
		file := filesdkmocks.NewMockJobFile(t)
		file.EXPECT().Remove(ctx, addProjectUID, addJobUID, api.RemoveRequest{FileURI: "a.json"}).
			Return(api.Result{SuccessCount: 1}, nil)

		s := service{jobFile: file, job: job, listFiles: stubList("a.json", "b.txt")}
		got, err := s.RunRemove(ctx, RemoveParams{ProjectID: addProjectUID, JobUIDOrName: addJobUID, FilePatterns: []string{"*.json"}})
		if err != nil {
			t.Fatalf("RunRemove: %v", err)
		}
		if got.Action != "removed" || got.SuccessCount != 1 || len(got.Files) != 1 || got.Files[0].FileURI != "a.json" {
			t.Errorf("unexpected output: %+v", got)
		}
	})
}
