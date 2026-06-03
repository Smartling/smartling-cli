package jobs

import (
	"context"
	"errors"
	"reflect"
	"testing"

	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	"github.com/Smartling/api-sdk-go/api/job"
	"github.com/stretchr/testify/mock"
)

func Test_service_RunProgress(t *testing.T) {
	const (
		validProjectUID = "PROJ87654321"
		validJobUID     = "aabbccdd1122"
		validJobName    = "Website Q1 2026"
	)
	apiErr := errors.New("api boom")

	tests := []struct {
		name    string
		setup   func(t *testing.T) *jobmocks.MockJob
		ctx     context.Context
		params  ProgressParams
		want    ProgressOutput
		wantErr bool
	}{
		{
			name:    "invalid: empty ProjectUID",
			setup:   func(t *testing.T) *jobmocks.MockJob { return jobmocks.NewMockJob(t) },
			ctx:     context.Background(),
			params:  ProgressParams{JobUIDOrName: validJobUID},
			wantErr: true,
		},
		{
			name:    "invalid: empty JobUIDOrName",
			setup:   func(t *testing.T) *jobmocks.MockJob { return jobmocks.NewMockJob(t) },
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID},
			wantErr: true,
		},
		{
			name: "UID lookup: GetJob succeeds, Progress returns output",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("GetJob", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobResponse{TranslationJobUID: validJobUID, JobName: validJobName}, nil)
				m.On("Progress", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobProgressResponse{
						TranslationJobUID: validJobUID,
						TotalWordCount:    100,
						PercentComplete:   42.5,
						JSON:              []byte(`{"percentComplete":42.5}`),
					}, nil)
				return m
			},
			ctx:    context.Background(),
			params: ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobUID},
			want: ProgressOutput{
				TranslationJobUID: validJobUID,
				TotalWordCount:    100,
				PercentComplete:   42.5,
				JSON:              []byte(`{"percentComplete":42.5}`),
			},
		},
		{
			name: "UID lookup: GetJob returns ErrNotFound, falls back to ListProjectJobs",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("GetJob", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobResponse{}, job.ErrNotFound)
				m.On("ListProjectJobs", mock.Anything, validProjectUID, job.ListProjectJobsParams{JobName: validJobUID}).
					Return(job.ListJobsResponse{Items: []job.JobSummary{
						{TranslationJobUID: "resolveduid01", JobName: validJobUID},
					}}, nil)
				m.On("Progress", mock.Anything, validProjectUID, "resolveduid01").
					Return(job.GetJobProgressResponse{
						TranslationJobUID: "resolveduid01",
						TotalWordCount:    7,
						PercentComplete:   10,
					}, nil)
				return m
			},
			ctx:    context.Background(),
			params: ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobUID},
			want: ProgressOutput{
				TranslationJobUID: "resolveduid01",
				TotalWordCount:    7,
				PercentComplete:   10,
			},
		},
		{
			name: "UID lookup: GetJob returns unexpected error",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("GetJob", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobResponse{}, apiErr)
				return m
			},
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobUID},
			wantErr: true,
		},
		{
			name: "name lookup: ListProjectJobs returns API error",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("ListProjectJobs", mock.Anything, validProjectUID, job.ListProjectJobsParams{JobName: validJobName}).
					Return(job.ListJobsResponse{}, apiErr)
				return m
			},
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobName},
			wantErr: true,
		},
		{
			name: "name lookup: ListProjectJobs returns empty list",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("ListProjectJobs", mock.Anything, validProjectUID, job.ListProjectJobsParams{JobName: validJobName}).
					Return(job.ListJobsResponse{Items: []job.JobSummary{}}, nil)
				return m
			},
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobName},
			wantErr: true,
		},
		{
			name: "name lookup: results have no matching name",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("ListProjectJobs", mock.Anything, validProjectUID, job.ListProjectJobsParams{JobName: validJobName}).
					Return(job.ListJobsResponse{Items: []job.JobSummary{
						{TranslationJobUID: "otheruid0001", JobName: "different name"},
					}}, nil)
				return m
			},
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobName},
			wantErr: true,
		},
		{
			name: "name lookup: matching name resolves to UID, Progress succeeds",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("ListProjectJobs", mock.Anything, validProjectUID, job.ListProjectJobsParams{JobName: validJobName}).
					Return(job.ListJobsResponse{Items: []job.JobSummary{
						{TranslationJobUID: "otheruid0001", JobName: "different name"},
						{TranslationJobUID: "matcheduid002", JobName: validJobName},
					}}, nil)
				m.On("Progress", mock.Anything, validProjectUID, "matcheduid002").
					Return(job.GetJobProgressResponse{
						TranslationJobUID: "matcheduid002",
						TotalWordCount:    250,
						PercentComplete:   75,
					}, nil)
				return m
			},
			ctx:    context.Background(),
			params: ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobName},
			want: ProgressOutput{
				TranslationJobUID: "matcheduid002",
				TotalWordCount:    250,
				PercentComplete:   75,
			},
		},
		{
			name: "Progress call fails",
			setup: func(t *testing.T) *jobmocks.MockJob {
				m := jobmocks.NewMockJob(t)
				m.On("GetJob", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobResponse{TranslationJobUID: validJobUID}, nil)
				m.On("Progress", mock.Anything, validProjectUID, validJobUID).
					Return(job.GetJobProgressResponse{}, apiErr)
				return m
			},
			ctx:     context.Background(),
			params:  ProgressParams{ProjectUID: validProjectUID, JobUIDOrName: validJobUID},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				job: tt.setup(t),
			}
			got, err := s.RunProgress(tt.ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunProgress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressParams_Validate(t *testing.T) {
	tests := []struct {
		name         string
		ProjectUID   string
		JobUIDOrName string
		wantErr      bool
	}{
		{
			name:         "valid",
			ProjectUID:   "PROJ87654321",
			JobUIDOrName: "aabbccdd1122",
			wantErr:      false,
		},
		{
			name:         "valid: JobUIDOrName as name",
			ProjectUID:   "PROJ87654321",
			JobUIDOrName: "Website Q1 2026",
			wantErr:      false,
		},
		{
			name:         "invalid: empty ProjectUID",
			ProjectUID:   "",
			JobUIDOrName: "aabbccdd1122",
			wantErr:      true,
		},
		{
			name:         "invalid: empty JobUIDOrName",
			ProjectUID:   "PROJ87654321",
			JobUIDOrName: "",
			wantErr:      true,
		},
		{
			name:         "invalid: all empty",
			ProjectUID:   "",
			JobUIDOrName: "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProgressParams{
				ProjectUID:   tt.ProjectUID,
				JobUIDOrName: tt.JobUIDOrName,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
