package joblist

import (
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// NewListCmd builds the `jobs list` command.
func NewListCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List translation jobs in a project or account.",
		Long: `List jobs within the configured project (default), within the account
(--account), or search jobs containing specific files or string hashcodes
(--file / --hashcode).`,
		Example: `
# List jobs in the current project

  smartling-cli jobs list

# Filter by name and status

  smartling-cli jobs list --name "Release" --status IN_PROGRESS

# List jobs across the account

  smartling-cli jobs list --account --with-priority

# Search jobs that contain a file

  smartling-cli jobs list --file path/to/a.json
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cnf, err := rootcmd.Config()
			if err != nil {
				return err
			}

			params, err := resolveParams(cmd, cnf.ProjectID, cnf.AccountID)
			if err != nil {
				return fmt.Errorf("failed to resolve list params: %w", err)
			}

			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	registerListFlags(listCmd)
	return listCmd
}

// registerListFlags adds all jobs-list flags to the command.
func registerListFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(accountFlag, false, "List jobs across the account instead of a single project.")
	cmd.Flags().String(nameFlag, "", "Filter by job name (maps to jobName).")
	cmd.Flags().String(numberFlag, "", "Filter by job number (project scope only).")
	cmd.Flags().StringArray(statusFlag, nil, "Filter by job status (repeatable; maps to translationJobStatus).")
	cmd.Flags().StringArray(uidFlag, nil, "Filter/search by translation job UID (repeatable).")
	cmd.Flags().StringArray(projectIDFlag, nil, "Filter by project ID (account scope only; repeatable).")
	cmd.Flags().StringArray(fileFlag, nil, "Search jobs containing this file URI (repeatable; uses the search endpoint).")
	cmd.Flags().StringArray(hashcodeFlag, nil, "Search jobs containing this string hashcode (repeatable; uses the search endpoint).")
	cmd.Flags().Bool(withPriorityFlag, false, "Include priority (account scope only).")
	cmd.Flags().String(sortByFlag, "", "Sort field.")
	cmd.Flags().String(sortDirectionFlag, "", "Sort direction (asc/desc).")
	cmd.Flags().Uint32(limitFlag, 0, "Maximum number of jobs to return.")
	cmd.Flags().Uint32(offsetFlag, 0, "Offset for pagination.")
}
