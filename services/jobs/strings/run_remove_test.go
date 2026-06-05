package jobstrings

import (
	"context"
	"errors"
	"testing"

	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"
	stringsdkmocks "github.com/Smartling/smartling-cli/services/jobs/strings/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/string"
)

func TestRunRemove(t *testing.T) {
	ctx := context.Background()
	const (
		projectUID = "test-project-id"
		jobUID     = "aabbccdd1122"
	)

	tests := []struct {
		name    string
		params  RemoveParams
		setup   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString)
		wantErr bool
		check   func(*testing.T, MutateOutput)
	}{
		{
			name:    "validation error — no hashcodes",
			params:  RemoveParams{ProjectID: projectUID, JobUIDOrName: jobUID},
			setup:   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString) {},
			wantErr: true,
		},
		{
			name: "resolves UID and removes strings",
			params: RemoveParams{
				ProjectID: projectUID, JobUIDOrName: jobUID,
				Hashcodes: []string{"h1"}, LocaleIDs: []string{"fr-FR"},
			},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().Remove(ctx, projectUID, jobUID, api.RemoveRequest{
					Hashcodes: []string{"h1"}, LocaleIDs: []string{"fr-FR"},
				}).Return(api.Result{SuccessCount: 1}, nil)
			},
			check: func(t *testing.T, got MutateOutput) {
				if got.Action != "removed" || got.TranslationJobUID != jobUID || len(got.Hashcodes) != 1 {
					t.Fatalf("unexpected output: %+v", got)
				}
				if got.SuccessCount != 1 {
					t.Errorf("SuccessCount = %d, want 1", got.SuccessCount)
				}
			},
		},
		{
			name:   "remove API error",
			params: RemoveParams{ProjectID: projectUID, JobUIDOrName: jobUID, Hashcodes: []string{"h1"}},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().Remove(ctx, projectUID, jobUID, api.RemoveRequest{Hashcodes: []string{"h1"}}).
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
			got, err := service{jobString: str, job: job}.RunRemove(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
