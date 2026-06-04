package locales

import (
	"context"
	"errors"
	"testing"

	localesdkmocks "github.com/Smartling/smartling-cli/services/jobs/locales/sdkmocks"
	jobsdkmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

func TestRunRemove(t *testing.T) {
	ctx := context.Background()
	const (
		projectUID = "test-project-id"
		jobUID     = "aabbccdd1122"
		localeID   = "fr-FR"
	)

	tests := []struct {
		name    string
		params  RemoveParams
		setup   func(*jobsdkmocks.MockJob, *localesdkmocks.MockJobLocale)
		wantErr bool
		check   func(*testing.T, Output)
	}{
		{
			name:    "validation error — empty job",
			params:  RemoveParams{ProjectID: projectUID, TargetLocaleID: localeID},
			setup:   func(*jobsdkmocks.MockJob, *localesdkmocks.MockJobLocale) {},
			wantErr: true,
		},
		{
			name:   "resolves UID and removes locale",
			params: RemoveParams{ProjectID: projectUID, JobUIDOrName: jobUID, TargetLocaleID: localeID},
			setup: func(j *jobsdkmocks.MockJob, l *localesdkmocks.MockJobLocale) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				l.EXPECT().Remove(ctx, projectUID, jobUID, localeID).Return(nil)
			},
			check: func(t *testing.T, got Output) {
				if got.Action != "removed" || got.TranslationJobUID != jobUID || got.TargetLocaleID != localeID {
					t.Fatalf("unexpected output: %+v", got)
				}
				if len(got.JSON) == 0 {
					t.Error("JSON should not be empty")
				}
			},
		},
		{
			name:   "remove API error",
			params: RemoveParams{ProjectID: projectUID, JobUIDOrName: jobUID, TargetLocaleID: localeID},
			setup: func(j *jobsdkmocks.MockJob, l *localesdkmocks.MockJobLocale) {
				j.EXPECT().GetJob(ctx, projectUID, jobUID).Return(jobapi.GetJobResponse{TranslationJobUID: jobUID}, nil)
				l.EXPECT().Remove(ctx, projectUID, jobUID, localeID).Return(errors.New("api error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := jobsdkmocks.NewMockJob(t)
			locale := localesdkmocks.NewMockJobLocale(t)
			tt.setup(job, locale)
			got, err := service{locale: locale, job: job}.RunRemove(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
