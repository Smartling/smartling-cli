package jobs

import (
	"context"
	"reflect"
	"testing"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	"github.com/Smartling/api-sdk-go/helpers/uid"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunList_ProjectScope(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("ListProjectJobs", mock.Anything, "proj-1", mock.MatchedBy(func(p jobapi.ListProjectJobsParams) bool {
		return p.JobName == "Release"
	})).Return(jobapi.ListJobsResponse{
		Items:      []jobapi.JobSummary{{TranslationJobUID: "u1", JobName: "Release", JobStatus: "IN_PROGRESS"}},
		TotalCount: 1,
	}, nil)

	s := NewService(m)
	out, err := s.RunList(context.Background(), ListParams{
		AccountUID: uid.AccountUID("test-account-uid"),
		ProjectUID: "proj-1",
		JobName:    "Release",
	})
	require.NoError(t, err)
	require.Len(t, out.Jobs, 1)
	require.Equal(t, "u1", out.Jobs[0].TranslationJobUID)
}

func TestRunList_AccountScope(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("ListAccountJobs", mock.Anything, "test-account-uid", mock.Anything).
		Return(jobapi.ListJobsResponse{Items: []jobapi.JobSummary{{TranslationJobUID: "a1"}}}, nil)

	s := NewService(m)
	out, err := s.RunList(context.Background(), ListParams{
		AccountUID: uid.AccountUID("test-account-uid"),
		ProjectUID: "proj-1",
		Account:    true,
	})
	require.NoError(t, err)
	require.Len(t, out.Jobs, 1)
	require.Equal(t, "a1", out.Jobs[0].TranslationJobUID)
}

func TestRunList_SearchScopeRejectsIncompatibleFlags(t *testing.T) {
	m := jobmocks.NewMockJob(t)

	s := NewService(m)
	_, err := s.RunList(context.Background(), ListParams{
		ProjectUID: "proj-1",
		FileURIs:   []string{"a.json"},
		JobStatus:  []string{"IN_PROGRESS"},
		Account:    true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--all-projects")
	require.Contains(t, err.Error(), "--status")
	m.AssertNotCalled(t, "SearchJobs", mock.Anything, mock.Anything, mock.Anything)
}

func TestRunList_SearchScope(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("SearchJobs", mock.Anything, "proj-1", mock.MatchedBy(func(r jobapi.SearchJobsRequest) bool {
		return len(r.FileURIs) == 1 && r.FileURIs[0] == "a.json"
	})).Return(jobapi.ListJobsResponse{Items: []jobapi.JobSummary{{TranslationJobUID: "s1"}}}, nil)

	s := NewService(m)
	out, err := s.RunList(context.Background(), ListParams{
		AccountUID: uid.AccountUID("test-account-uid"),
		ProjectUID: "proj-1",
		FileURIs:   []string{"a.json"},
	})
	require.NoError(t, err)
	require.Len(t, out.Jobs, 1)
	require.Equal(t, "s1", out.Jobs[0].TranslationJobUID)
}

func TestRunList_TruncationNote(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("ListProjectJobs", mock.Anything, "proj-1", mock.Anything).
		Return(jobapi.ListJobsResponse{
			Items:      []jobapi.JobSummary{{TranslationJobUID: "u1"}, {TranslationJobUID: "u2"}},
			TotalCount: 10,
		}, nil)

	s := NewService(m)
	out, err := s.RunList(context.Background(), ListParams{
		ProjectUID: "proj-1",
		Limit:      2,
		Offset:     4,
	})
	require.NoError(t, err)
	require.Equal(t, 10, out.TotalCount)
	require.Equal(t, uint32(4), out.Offset)

	lines := out.SimpleLines()
	require.Contains(t, lines[len(lines)-1], "Showing 2 of 10 jobs")
	require.Contains(t, lines[len(lines)-1], "--offset 6")
}

func TestRunList_NoTruncationWhenComplete(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("ListProjectJobs", mock.Anything, "proj-1", mock.Anything).
		Return(jobapi.ListJobsResponse{
			Items:      []jobapi.JobSummary{{TranslationJobUID: "u1"}, {TranslationJobUID: "u2"}},
			TotalCount: 2,
		}, nil)

	s := NewService(m)
	out, err := s.RunList(context.Background(), ListParams{ProjectUID: "proj-1"})
	require.NoError(t, err)
	for _, line := range out.SimpleLines() {
		require.NotContains(t, line, "Showing")
	}
}

func TestListParams_searchConflicts(t *testing.T) {
	tests := []struct {
		name   string
		params ListParams
		want   []string
	}{
		{
			name:   "no conflicts",
			params: ListParams{FileURIs: []string{"a.json"}},
			want:   nil,
		},
		{
			name:   "all-projects",
			params: ListParams{Account: true},
			want:   []string{"--all-projects"},
		},
		{
			name:   "name",
			params: ListParams{JobName: "Release"},
			want:   []string{"--name"},
		},
		{
			name:   "number",
			params: ListParams{JobNumber: "42"},
			want:   []string{"--number"},
		},
		{
			name:   "status",
			params: ListParams{JobStatus: []string{"IN_PROGRESS"}},
			want:   []string{"--status"},
		},
		{
			name:   "project-id",
			params: ListParams{ProjectIDs: []string{"p1"}},
			want:   []string{"--project-id"},
		},
		{
			name:   "with-priority",
			params: ListParams{WithPriority: true},
			want:   []string{"--with-priority"},
		},
		{
			name:   "sort-by",
			params: ListParams{SortBy: "name"},
			want:   []string{"--sort-by"},
		},
		{
			name:   "sort-direction",
			params: ListParams{SortDirection: "asc"},
			want:   []string{"--sort-direction"},
		},
		{
			name:   "limit is not a conflict (pagination ignored by search)",
			params: ListParams{Limit: 10},
			want:   nil,
		},
		{
			name:   "offset is not a conflict (pagination ignored by search)",
			params: ListParams{Offset: 5},
			want:   nil,
		},
		{
			name: "all conflicts in declared order",
			params: ListParams{
				Account:       true,
				JobName:       "Release",
				JobNumber:     "42",
				JobStatus:     []string{"IN_PROGRESS"},
				ProjectIDs:    []string{"p1"},
				WithPriority:  true,
				SortBy:        "name",
				SortDirection: "asc",
				Limit:         10,
				Offset:        5,
			},
			want: []string{
				"--all-projects",
				"--name",
				"--number",
				"--status",
				"--project-id",
				"--with-priority",
				"--sort-by",
				"--sort-direction",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.searchConflicts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchConflicts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  ListParams
		wantErr bool
	}{
		{
			name:    "project scope ok",
			params:  ListParams{ProjectUID: "proj-1"},
			wantErr: false,
		},
		{
			name:    "project scope missing project uid",
			params:  ListParams{},
			wantErr: true,
		},
		{
			name:    "account scope ok",
			params:  ListParams{Account: true, AccountUID: uid.AccountUID("test-account-uid")},
			wantErr: false,
		},
		{
			name:    "account scope missing account uid",
			params:  ListParams{Account: true},
			wantErr: true,
		},
		{
			name:    "search scope ok",
			params:  ListParams{ProjectUID: "proj-1", FileURIs: []string{"a.json"}},
			wantErr: false,
		},
		{
			name:    "search scope by hashcode ok",
			params:  ListParams{ProjectUID: "proj-1", Hashcodes: []string{"h1"}},
			wantErr: false,
		},
		{
			name:    "search scope missing project uid",
			params:  ListParams{FileURIs: []string{"a.json"}},
			wantErr: true,
		},
		{
			name: "search scope with conflicting flag",
			params: ListParams{
				ProjectUID: "proj-1",
				FileURIs:   []string{"a.json"},
				JobStatus:  []string{"IN_PROGRESS"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.params.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
