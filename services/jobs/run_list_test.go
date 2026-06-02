package jobs

import (
	"context"
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
