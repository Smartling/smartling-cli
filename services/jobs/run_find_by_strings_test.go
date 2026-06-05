package jobs

import (
	"context"
	"testing"
	"time"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	jobmocks "github.com/Smartling/smartling-cli/services/jobs/sdkmocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunFindByStrings_FlattensMatches(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("FindJobsByStrings", mock.Anything, "proj-1", mock.MatchedBy(func(r jobapi.FindJobsByStringsRequest) bool {
		return len(r.Hashcodes) == 1 && r.Hashcodes[0] == "h1" &&
			len(r.LocaleIDs) == 1 && r.LocaleIDs[0] == "fr-FR"
	})).Return(jobapi.FindJobsByStringsResponse{
		TotalCount: 1,
		Items: []jobapi.JobWithStrings{
			{
				TranslationJobUID: "u1",
				JobName:           "Release",
				HashcodesByLocale: []jobapi.JobHashcodesByLocale{
					{LocaleID: "fr-FR", Hashcodes: []string{"h1"}},
					{LocaleID: "de-DE", Hashcodes: []string{"h1"}},
				},
			},
		},
	}, nil)

	s := NewService(m)
	out, err := s.RunFindByStrings(context.Background(), FindByStringsParams{
		ProjectUID: "proj-1",
		Hashcodes:  []string{"h1"},
		LocaleIDs:  []string{"fr-FR"},
	})
	require.NoError(t, err)
	require.Len(t, out.Matches, 2)
	require.Equal(t, FindByStringsMatch{
		Hashcode: "h1", LocaleID: "fr-FR",
		TranslationJobUID: "u1", JobName: "Release",
	}, out.Matches[0])
	require.Equal(t, "de-DE", out.Matches[1].LocaleID)
	require.Equal(t, "u1", out.Matches[1].TranslationJobUID)
}

func TestRunFindByStrings_NoMatches(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	m.On("FindJobsByStrings", mock.Anything, "proj-1", mock.Anything).
		Return(jobapi.FindJobsByStringsResponse{}, nil)

	s := NewService(m)
	out, err := s.RunFindByStrings(context.Background(), FindByStringsParams{
		ProjectUID: "proj-1",
		Hashcodes:  []string{"h1"},
	})
	require.NoError(t, err)
	require.Empty(t, out.Matches)
	require.Equal(t, []string{"No jobs found for the given strings."}, out.SimpleLines())
}

func TestFindByStringsParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  FindByStringsParams
		wantErr bool
	}{
		{
			name:    "ok",
			params:  FindByStringsParams{ProjectUID: "proj-1", Hashcodes: []string{"h1"}},
			wantErr: false,
		},
		{
			name:    "missing project uid",
			params:  FindByStringsParams{Hashcodes: []string{"h1"}},
			wantErr: true,
		},
		{
			name:    "missing hashcodes",
			params:  FindByStringsParams{ProjectUID: "proj-1"},
			wantErr: true,
		},
		{
			name:    "at record limit (hashcodes only)",
			params:  FindByStringsParams{ProjectUID: "proj-1", Hashcodes: makeStrings(MaxFindByStringsRecords)},
			wantErr: false,
		},
		{
			name:    "over record limit by hashcodes alone",
			params:  FindByStringsParams{ProjectUID: "proj-1", Hashcodes: makeStrings(MaxFindByStringsRecords + 1)},
			wantErr: true,
		},
		{
			name: "over record limit by hashcodes×locales",
			params: FindByStringsParams{
				ProjectUID: "proj-1",
				Hashcodes:  makeStrings(2001),
				LocaleIDs:  makeStrings(10),
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

func TestRunFindByStrings_RendersDueDateAndCount(t *testing.T) {
	due := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	m := jobmocks.NewMockJob(t)
	m.On("FindJobsByStrings", mock.Anything, "proj-1", mock.Anything).
		Return(jobapi.FindJobsByStringsResponse{
			TotalCount: 2,
			Items: []jobapi.JobWithStrings{
				{
					TranslationJobUID: "u1",
					JobName:           "Has due date",
					DueDate:           due,
					HashcodesByLocale: []jobapi.JobHashcodesByLocale{
						{LocaleID: "fr-FR", Hashcodes: []string{"h1"}},
					},
				},
				{
					TranslationJobUID: "u2",
					JobName:           "No due date",
					// zero DueDate -> empty string
					HashcodesByLocale: []jobapi.JobHashcodesByLocale{
						{LocaleID: "de-DE", Hashcodes: []string{"h1"}},
					},
				},
			},
		}, nil)

	s := NewService(m)
	out, err := s.RunFindByStrings(context.Background(), FindByStringsParams{
		ProjectUID: "proj-1",
		Hashcodes:  []string{"h1"},
	})
	require.NoError(t, err)
	require.Equal(t, 2, out.TotalCount)
	require.Len(t, out.Matches, 2)

	// DueDate is RFC3339 for the dated job, empty for the undated one.
	require.Equal(t, "2026-01-02T03:04:05Z", out.Matches[0].DueDate)
	require.Equal(t, "", out.Matches[1].DueDate)

	// TableData mirrors the matches with a fixed header row.
	headers, rows := out.TableData()
	require.Equal(t, []string{"HASHCODE", "LOCALE", "TRANSLATION JOB UID", "JOB NAME", "DUE DATE"}, headers)
	require.Equal(t, [][]string{
		{"h1", "fr-FR", "u1", "Has due date", "2026-01-02T03:04:05Z"},
		{"h1", "de-DE", "u2", "No due date", ""},
	}, rows)

	// SimpleLines ends with the job-count summary.
	lines := out.SimpleLines()
	require.Equal(t, "2 job(s) matched.", lines[len(lines)-1])
}

func makeStrings(n int) []string {
	s := make([]string, n)
	for i := range s {
		s[i] = "h"
	}
	return s
}

func TestRunFindByStrings_InvalidParams(t *testing.T) {
	m := jobmocks.NewMockJob(t)
	s := NewService(m)
	_, err := s.RunFindByStrings(context.Background(), FindByStringsParams{Hashcodes: []string{"h1"}})
	require.Error(t, err)
	m.AssertNotCalled(t, "FindJobsByStrings", mock.Anything, mock.Anything, mock.Anything)
}
