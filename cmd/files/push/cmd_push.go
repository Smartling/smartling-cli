package push

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/help"

	"github.com/spf13/cobra"
)

// NewPushCmd creates a new command to upload files to the Smartling platform.
func NewPushCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	var (
		authorize  bool
		locales    []string
		branch     string
		fileType   string
		directory  string
		directives []string
		job        string
	)

	pushCmd := &cobra.Command{
		Use:     "push <file> <uri> --job <job name> [--authorize] [--locale <locale>]",
		Aliases: []string{"upload"},
		Short:   "Creates job and uploads specified file into this job.",
		Long: `smartling-cli files push <file> [<uri>] --job <job name> [--authorize] [--locale <locale>] [--type <type>] [--branch (@auto|<branch name>)] [--directory <work dir>] [--directive <smartling directive>]

Creates a new job (or reuses existing) in Smartling TMS and uploads designated
file(s) for translation.

One or more files can be pushed.

When pushing single file, <uri> can be specified to override local path.
When pushing multiple files, they will be uploaded using local path as URI.
If no file specified in command line, config file will be used to lookup
for file masks to push.

Use --job option to specify job name or job UID. If job name is not specified,
then the "CLI uploads" name will be used.
You can use the same job name for multiple CLI calls. The same job will be used
in this case (CLI searches by the job name). If the job with the same name exists,
but it has state Canceled or Closed, then a new job will be created with timestamp suffix.

To authorize the job after uploading all files, use --authorize option.

To specify locales for the files in the job, use one or more --locale options.
If no locales specified, then all project locales will be added to all uploaded files.

To prepend prefix to all target URIs, use --branch option. Special
value "@auto" can be used to tell the tool to use the current git
branch name as value for --branch option.

File type will be deduced from file extension. If file extension is unknown,
type should be specified manually by using --type option. That option also
can be used to override detected file type.

` + "`<file>` " + help.GlobPattern + ` 

` + help.AuthenticationOptions,
		Example: `
# Upload files to a translation job

  smartling-cli files push my-file.txt --job "Website Update" --authorize

# Upload multiple files using pattern matching with the command alias ‘upload’

  smartling-cli files upload "src/**/*.json" --job "App Localization"

# Manual branch naming

  smartling-cli files push "**/*.txt" --branch "feature-branch"

# Automatic Git branch detection

  smartling-cli files push "**/*.txt" --branch "@auto"

# All JSON files in subdirectories

  smartling-cli files push "**/*.json"

# Specific file types

  smartling-cli files push "**/*.{json,xml,properties}"

# Files matching naming convention with the command alias 'upload' 

  smartling-cli files upload "**/messages_*.properties"

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var (
				file string
				uri  string
			)
			if len(args) > 0 {
				file = args[0]
			}
			if len(args) > 1 {
				uri = args[1]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				return err
			}

			directives, err := helpers.MKeyValueToMap(directives)
			if err != nil {
				return err
			}
			p := files.PushParams{
				URI:         uri,
				File:        file,
				Branch:      branch,
				Locales:     locales,
				Authorize:   authorize,
				Directory:   directory,
				FileType:    fileType,
				Directives:  directives,
				JobIDOrName: job,
			}

			return s.RunPush(ctx, p)
		},
	}

	pushCmd.Flags().BoolVarP(&authorize, "authorize", "z", false, `Automatically authorize the job with file(s) and specified locales.
If the flag is not specified, the job remains unauthorized.`)
	pushCmd.Flags().StringArrayVarP(&locales, "locale", "l", []string{}, `<locale code>
Add file(s) to the job for the specified locale only.
If the flag is not specified, then all project locales will be added to the job.
Can be specified several times: --locale fr --locale de -l es`)
	pushCmd.Flags().StringVarP(&branch, "branch", "b", "", `<branch>
Prepend specified prefix to target file URI.`)
	pushCmd.Flags().StringVarP(&fileType, "type", "t", "", `<type>
Override automatically detected file type.`)
	pushCmd.Flags().StringArrayVarP(&directives, "directive", "r", []string{}, `Specify one or more directives to use in push request.`)
	pushCmd.Flags().StringVarP(&directory, "directory", "d", ".", `Specified directory.`)
	pushCmd.Flags().StringVarP(&job, "job", "j", "", `<job name>
Provide a name for the Smartling translation job or job UID.
All files will be uploaded into this job.
If the flag is not specified then the "CLI uploads" name will be used.`)

	return pushCmd
}
