package joblist

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestResolveParams_DefaultProjectScope(t *testing.T) {
	cmd := newTestCmd()
	require.NoError(t, cmd.Flags().Set("name", "Release"))

	params, err := resolveParams(cmd, "proj-1", "test-account-uid")
	require.NoError(t, err)
	require.False(t, params.Account)
	require.Equal(t, "Release", params.JobName)
	require.Equal(t, "proj-1", params.ProjectUID)
}

func TestResolveParams_AccountScope(t *testing.T) {
	cmd := newTestCmd()
	require.NoError(t, cmd.Flags().Set("account", "true"))

	params, err := resolveParams(cmd, "proj-1", "test-account-uid")
	require.NoError(t, err)
	require.True(t, params.Account)
}

// newTestCmd builds a command carrying the same flags as the real one so
// resolveParams can read them.
func newTestCmd() *cobra.Command {
	c := &cobra.Command{Use: "list", RunE: func(*cobra.Command, []string) error { return nil }}
	registerListFlags(c)
	return c
}
