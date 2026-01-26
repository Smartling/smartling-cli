package progress

import (
	"errors"
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/spf13/cobra"
)

// NewProgressCmd returns new progress command
func NewProgressCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	progressCmd := &cobra.Command{
		Use:   "progress <translationJobUid|translationJobName>",
		Short: "Track translation progress for a specific job.",
		Long: `smartling-cli jobs progress <translationJobUid|translationJobName> [--output json]

Retrieves real-time translation progress metrics for a specific translation job.
This command is essential for monitoring active translations, estimating completion times,
and tracking workflow progress across multiple locales.

Progress information includes total word counts, completion percentages, and detailed
per-locale breakdowns showing how content moves through each translation workflow step
(awaiting authorization, in translation, completed, etc.).

The command accepts either:
  • Translation Job UID: 12-character alphanumeric identifier (e.g., aabbccdd1122)
  • Translation Job Name: Human-readable name assigned when creating the job

If multiple jobs share the same name, the most recent active job (not Canceled or Closed)
will be selected.

Output Formats:

  --output simple (default)
    Displays key progress metrics in human-readable format:
      - Total word count across all locales
      - Overall completion percentage
    Best for: Quick status checks, manual monitoring, terminal viewing

  --output json
    Returns the complete API response as JSON, including:
      - Per-locale progress breakdowns
      - Workflow step details (authorized, awaiting, completed, etc.)
      - String and word counts at each workflow stage
      - Target locale descriptions
    Best for: Automation scripts, CI/CD pipelines, custom reporting tools

Use Cases:
  • Monitor active translation projects to estimate delivery times
  • Track progress before authorizing next workflow steps
  • Build automated alerts when translations reach completion thresholds
  • Generate custom progress reports for stakeholders
  • Integrate with CI/CD pipelines to gate deployments on translation completion

Project Configuration:
  Project ID must be configured in smartling.yml or specified via --project flag.
  Account ID can be configured in smartling.yml or specified via --account flag.

Authentication is required via user_id and secret in smartling.yml or environment variables.

Available options:` + help.AuthenticationOptions,
		Example: `
# Check progress using job name

  smartling-cli jobs progress "Website Q1 2026"

# Check progress using job UID

  smartling-cli jobs progress aabbccdd1122

# Get detailed JSON output for automation

  smartling-cli jobs progress "Mobile App Release" --output json

# Use with specific project

  smartling-cli jobs progress aabbccdd1122 --project 9876543210

# Parse JSON output in scripts (example: check if job is 100% complete)

  PROGRESS=$(smartling-cli jobs progress my-job --output json | jq '.percentComplete')
  if [ "$PROGRESS" -eq 100 ]; then
    echo "Translation complete!"
  fi

# Monitor progress for CI/CD gate

  smartling-cli jobs progress "Release v2.0" --output json | \
    jq -e '.percentComplete >= 95' && echo "Ready for deployment"

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 1 {
				return clierror.UIError{
					Operation:   "check args",
					Err:         errors.New("wrong argument quantity"),
					Description: fmt.Sprintf("expected one argument, got: %d", len(args)),
				}
			}
			idOrName := args[0]

			cnf, err := rootcmd.Config()
			if err != nil {
				return err
			}

			accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
			if err != nil {
				return err
			}

			params := srv.ProgressParams{
				AccountUID:   accountUID,
				ProjectUID:   cnf.ProjectID,
				JobUIDOrName: idOrName,
			}
			format, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	return progressCmd
}
