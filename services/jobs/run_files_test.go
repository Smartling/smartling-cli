package jobs

import (
	"context"
	"testing"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunFiles_ResolvesUIDAndReturnsFiles(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("GetJob", mock.Anything, "proj-1", "aabbccdd1122").
		Return(jobapi.GetJobResponse{TranslationJobUID: "aabbccdd1122", JobName: "My Job"}, nil)
	m.On("ListFiles", mock.Anything, "proj-1", "aabbccdd1122", uint32(500), uint32(0)).
		Return(jobapi.ListJobFilesResponse{
			TotalCount: 1,
			Items:      []jobapi.JobFile{{FileURI: "/a.json", LocaleIDs: []string{"fr-FR"}}},
		}, nil)

	s := NewService(m)
	out, err := s.RunFiles(context.Background(), FilesParams{
		ProjectUID:   "proj-1",
		JobUIDOrName: "aabbccdd1122",
		Limit:        500,
		Offset:       0,
	})
	require.NoError(t, err)
	require.Len(t, out.Files, 1)
	require.Equal(t, "/a.json", out.Files[0].FileURI)
	require.Equal(t, 1, out.TotalCount)
}
