package jobsfiles

import (
	"context"
	"testing"

	filesdkmocks "github.com/Smartling/smartling-cli/services/jobs/files/sdkmocks"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	jobfile "github.com/Smartling/api-sdk-go/api/job/file"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunList_ResolvesUIDAndReturnsFiles(t *testing.T) {
	job := jobmocks.NewMockJob(t)
	job.On("GetJob", mock.Anything, "proj-1", "aabbccdd1122").
		Return(jobapi.GetJobResponse{TranslationJobUID: "aabbccdd1122", JobName: "My Job"}, nil)

	file := filesdkmocks.NewMockJobFile(t)
	file.On("List", mock.Anything, "proj-1", "aabbccdd1122", uint32(500), uint32(0)).
		Return(jobfile.ListResponse{
			TotalCount: 1,
			Items:      []jobfile.File{{FileURI: "/a.json", LocaleIDs: []string{"fr-FR"}}},
		}, nil)

	s := service{job: job, jobFile: file}
	out, err := s.RunList(context.Background(), ListParams{
		ProjectID:    "proj-1",
		JobUIDOrName: "aabbccdd1122",
		Limit:        500,
		Offset:       0,
	})
	require.NoError(t, err)
	require.Len(t, out.Files, 1)
	require.Equal(t, "/a.json", out.Files[0].FileURI)
	require.Equal(t, 1, out.TotalCount)
}
