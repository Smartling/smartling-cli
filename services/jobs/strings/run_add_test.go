package jobstrings

import (
	"context"
	"errors"
	"strings"
	"testing"

	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"
	stringsdkmocks "github.com/Smartling/smartling-cli/services/jobs/strings/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/string"
)

func TestRunAdd(t *testing.T) {
	ctx := context.Background()
	const (
		projectUID = "test-project-id"
		jobUID     = "aabbccdd1122"
	)

	tests := []struct {
		name    string
		params  AddParams
		setup   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString)
		wantErr bool
		check   func(*testing.T, MutateOutput)
	}{
		{
			name:    "validation error — no hashcodes",
			params:  AddParams{ProjectID: projectUID, JobUIDOrName: jobUID},
			setup:   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString) {},
			wantErr: true,
		},
		{
			name: "resolves UID and adds strings",
			params: AddParams{
				ProjectID: projectUID, JobUIDOrName: jobUID,
				Hashcodes: []string{"h1", "h2"}, TargetLocaleIDs: []string{"fr-FR"}, MoveEnabled: true,
			},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().Add(ctx, projectUID, jobUID, api.AddRequest{
					Hashcodes: []string{"h1", "h2"}, TargetLocaleIDs: []string{"fr-FR"}, MoveEnabled: true,
				}).Return(api.Result{SuccessCount: 2}, nil)
			},
			check: func(t *testing.T, got MutateOutput) {
				if got.Action != "added" || got.TranslationJobUID != jobUID || len(got.Hashcodes) != 2 {
					t.Fatalf("unexpected output: %+v", got)
				}
				if got.SuccessCount != 2 {
					t.Errorf("SuccessCount = %d, want 2", got.SuccessCount)
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name:   "nonexistent hashcodes report zero affected with a hint",
			params: AddParams{ProjectID: projectUID, JobUIDOrName: jobUID, Hashcodes: []string{"bad1", "bad2"}},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().Add(ctx, projectUID, jobUID, api.AddRequest{Hashcodes: []string{"bad1", "bad2"}}).
					Return(api.Result{SuccessCount: 0, FailCount: 0}, nil)
			},
			check: func(t *testing.T, got MutateOutput) {
				if got.SuccessCount != 0 || got.FailCount != 0 {
					t.Fatalf("counts = %d/%d, want 0/0", got.SuccessCount, got.FailCount)
				}
				lines := got.SimpleLines()
				if len(lines) != 2 || !strings.Contains(lines[1], "Verify the hashcodes") {
					t.Errorf("SimpleLines = %v, want a zero-affected hint", lines)
				}
			},
		},
		{
			name:   "job not found by name",
			params: AddParams{ProjectID: projectUID, JobUIDOrName: "No Such Job", Hashcodes: []string{"h1"}},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().ListProjectJobs(ctx, projectUID, jobapi.ListProjectJobsParams{JobName: "No Such Job"}).
					Return(jobapi.ListJobsResponse{}, nil)
			},
			wantErr: true,
		},
		{
			name:   "add API error",
			params: AddParams{ProjectID: projectUID, JobUIDOrName: jobUID, Hashcodes: []string{"h1"}},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().Add(ctx, projectUID, jobUID, api.AddRequest{Hashcodes: []string{"h1"}}).
					Return(api.Result{}, errors.New("api error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := jobsdkmocks.NewMockJob(t)
			str := stringsdkmocks.NewMockJobString(t)
			tt.setup(job, str)
			got, err := service{jobString: str, job: job}.RunAdd(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
