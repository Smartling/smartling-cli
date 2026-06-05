package jobstrings

import (
	"context"
	"strings"
	"testing"

	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"
	stringsdkmocks "github.com/Smartling/smartling-cli/services/jobs/strings/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	api "github.com/Smartling/api-sdk-go/api/job/string"
)

func TestRunList(t *testing.T) {
	ctx := context.Background()
	const (
		projectUID = "test-project-id"
		jobUID     = "aabbccdd1122"
	)

	tests := []struct {
		name    string
		params  ListParams
		setup   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString)
		wantErr bool
		check   func(*testing.T, ListOutput)
	}{
		{
			name:    "validation error — empty job",
			params:  ListParams{ProjectID: projectUID},
			setup:   func(*jobsdkmocks.MockJob, *stringsdkmocks.MockJobString) {},
			wantErr: true,
		},
		{
			name:   "resolves UID and lists strings with filter",
			params: ListParams{ProjectID: projectUID, JobUIDOrName: jobUID, TargetLocaleID: "fr-FR", Limit: 10},
			setup: func(j *jobsdkmocks.MockJob, s *stringsdkmocks.MockJobString) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				s.EXPECT().List(ctx, projectUID, jobUID, api.ListParams{TargetLocaleID: "fr-FR", Limit: 10}).
					Return(api.ListResponse{
						TotalCount: 2,
						Items: []api.StringHashcode{
							{TargetLocaleID: "fr-FR", Hashcode: "h1"},
							{TargetLocaleID: "fr-FR", Hashcode: "h2"},
						},
					}, nil)
			},
			check: func(t *testing.T, got ListOutput) {
				if got.TotalCount != 2 || len(got.Items) != 2 {
					t.Fatalf("unexpected output: %+v", got)
				}
				if got.Items[0].Hashcode != "h1" || got.Items[1].TargetLocaleID != "fr-FR" {
					t.Errorf("unexpected items: %+v", got.Items)
				}
				lines := got.SimpleLines()
				if !strings.Contains(lines[len(lines)-1], "Showing 2 of 2") {
					t.Errorf("SimpleLines = %v, want trailing total line", lines)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := jobsdkmocks.NewMockJob(t)
			str := stringsdkmocks.NewMockJobString(t)
			tt.setup(job, str)
			got, err := service{jobString: str, job: job}.RunList(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
