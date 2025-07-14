package push

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
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
		Use:   "push <file> <uri>",
		Short: "Uploads specified file into Smartling platform.",
		Long: `smartling-cli files push <file> [<uri>] [--type <type>] [--branch (@auto|<branch name>)] [--authorize|--locale <locale>] [--directory <work dir>] [--directive <smartling directive>]

Uploads files designated for translation.

One or several files can be pushed.

When pushing single file, <uri> can be specified to override local path.
When pushing multiple files, they will be uploaded using local path as URI.
If no file specified in command line, config file will be used to lookup
for file masks to push.

To authorize all locales, use --authorize option.

To authorize only specific locales, use one or more --locale.

To prepend prefix to all target URIs, use --branch option. Special
value "@auto" can be used to tell that tool should try to took current git
branch name as value for --branch option.

File type will be deduced from file extension. If file extension is unknown,
type should be specified manually by using --type option. That option also
can be used to override detected file type.

<file> ` + help.GlobPattern + `


Available options:
  -p --project <project>
    Specify project to use.

  --authorize
    Authorize all available locales. Incompatible with --locale option.

  --locale <locale>
    Authorize speicified locale only. Can be specified several times.
    Incompatible with --authorize option.

  --branch <branch>
    Prepend specified prefix to target file URI.

  --type <type>
    Override automatically detected file type.
` + help.AuthenticationOptions,
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

	pushCmd.Flags().BoolVarP(&authorize, "authorize", "z", false, `Automatically authorize all locales in specified file. Incompatible with -l option.`)
	pushCmd.Flags().StringArrayVarP(&locales, "locale", "l", []string{}, `Authorize only specified locales.`)
	pushCmd.Flags().StringVarP(&branch, "branch", "b", "", `Prepend specified text to the file uri.`)
	pushCmd.Flags().StringVarP(&fileType, "type", "t", "", `Specifies file type which will be used instead of automatically deduced from extension.`)
	pushCmd.Flags().StringArrayVarP(&directives, "directive", "r", []string{}, `Specifies one or more directives to use in push request.`)
	pushCmd.Flags().StringVarP(&directory, "directory", "d", ".", `Specified directory.`)
	pushCmd.Flags().StringVarP(&job, "job", "j", "", `Enter a name for the Smartling translation job or job UID. All files will be uploaded into this job.`)

	return pushCmd
}
