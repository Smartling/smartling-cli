## smartling-cli jobs progress

Track translation progress for a specific job.

### Synopsis

smartling-cli jobs progress <translationJobUid|translationJobName> [--output json]

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

Available options:
  --user <user>
    Specify user ID for authentication.

  --secret <secret>
    Specify secret token for authentication.

  -a --account <account>
    Specify account ID.


```
smartling-cli jobs progress <translationJobUid|translationJobName> [flags]
```

### Examples

```

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


```

### Options

```
  -h, --help   help for progress
```

### Options inherited from parent commands

```
  -a, --account string               Account ID to operate on.
                                     This option overrides config value "account_id".
  -c, --config string                Config file in YAML format.
                                     By default CLI will look for file named
                                     "smartling.yml" in current directory and in all
                                     intermediate parents, emulating git behavior.
  -k, --insecure                     Skip HTTPS certificate validation.
      --operation-directory string   Sets directory to operate on, usually, to store or to
                                     read files.  Depends on command. (default ".")
      --output string                Output format: json, simple (default "simple")
  -p, --project string               Project ID to operate on.
                                     This option overrides config value "project_id".
      --proxy string                 Use specified URL as proxy server.
      --secret string                Token Secret which will be used for authentication.
                                     This option overrides config value "secret".
      --smartling-url string         Specify base Smartling URL, merely for testing
                                     purposes.
      --threads uint32               If command can be executed concurrently, it will be
                                     executed for at most <number> of threads. (default 4)
      --user string                  User ID which will be used for authentication.
                                     This option overrides config value "user_id".
  -v, --verbose count                Verbose logging
```

### SEE ALSO

* [smartling-cli jobs](smartling-cli_jobs.md)	 - Manage translation jobs and monitor their progress.

###### Auto generated by spf13/cobra on 23-Jan-2026
