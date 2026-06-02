package jobs

import (
	"context"
	"testing"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunView_ResolvesUIDAndReturnsDetails(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	// 12-char lowercase-alnum input is treated as a UID by the resolver.
	m.On("GetJob", mock.Anything, "proj-1", "aabbccdd1122").
		Return(jobapi.GetJobResponse{
			TranslationJobUID: "aabbccdd1122",
			JobName:           "My Job",
			JobStatus:         "IN_PROGRESS",
			JobNumber:         "SMTL-7",
		}, nil)

	s := NewService(m)
	out, err := s.RunView(context.Background(), ViewParams{
		ProjectUID:   "proj-1",
		JobUIDOrName: "aabbccdd1122",
	})
	require.NoError(t, err)
	require.Equal(t, "My Job", out.JobName)
	require.Equal(t, "IN_PROGRESS", out.JobStatus)
}
