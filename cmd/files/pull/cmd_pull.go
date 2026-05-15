package pull

import (
	"os"
	"strconv"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

const threadsFlag = "threads"

var (
	uri        string
	jobUID     string
	all        bool
	source     bool
	resume     bool
	dryRun     bool
	progress   string
	retrieve   string
	directory  string
	formatPath string
	locales    []string
	threads    uint32
)

// NewPullCmd creates a new command to pull files.
func NewPullCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	pullCmd := &cobra.Command{
		Use:     "pull <uri>",
		Aliases: []string{"download"},
		Short:   "Pulls specified files from server.",
		Long: `smartling-cli files pull — downloads translated files from project.

Downloads files from specified project into local directory.

It's possible to download only specific files by file mask, to download source
files with translations, to download file to specific directory or to download
specific locales only.

If special value of "-" is specified as <uri>, then program will expect
to read files list from stdin:

  cat files.txt | smartling-cli files pull -

` + "`<uri>` " + help.GlobPattern + ` 

If --locale flag is not specified, all available locales are downloaded. To
see available locales, use "status" command.

To download files into subdirectory, use --directory option and specify
directory name you want to download into.

To download source file as well as translated files specify --source option.

Files will be downloaded and stored under names used while upload (e.g. File
URI). While downloading translated file suffix "_<locale>" will be appended to
file name before extension. To override file format name, use --format option.
` + help.FormatOption + `
Following variables are available:

  > .FileURI — full file URI in Smartling system;
  > .Locale — locale ID for translated file and empty for source file;
  > .JobUID — translation job UID, when --job-uid is set (otherwise empty);


Available options:
  -p --project <project>
    Specify project to use.

  --source
    Download source files along with translated files.

  --all
    Download all translated files. Required if no file pattern is specified.

  —d ——directory <dir>
    Download files into specified directory.

  --format <format>
    Specify format for download file name.

  --progress <percents>
    Specify minimum of translation progress in percents.
	By default that filter does not apply.

  --retrieve <type>
    Retrieval type according to API specs:
    > pending — returns any translations, including non-published ones);
    > published — returns only published translations;
    > pseudo — returns modified version of original text with certain
               characters transformed;
    > contextMatchingInstrumented — to use with Chrome Context Capture;

  --job-uid <jobUid>
    Download every file × target-locale pair from a Smartling job.
    Combines with the <uri> positional argument: when both are set, the
    URI is treated as a glob filter applied to the job's file list.
    Combines with --locale: the requested locales are intersected with the job's target locales.
    When --format is not set, defaults to <jobUid>/<locale>/<fileUri>.

  --resume
    Skip files that already exist on disk. Useful for re-running a large
    pull after a failure.

  --dry-run
    Print the file × locale matrix that would be downloaded, then exit 0.
    Does not call GetFileStatus, so --progress filtering is not applied.
` + help.AuthenticationOptions,
		Example: `
# Pull translated files

  smartling-cli files pull "**/*.json" --locale fr-FR --locale de-DE

# Use the alias 'download' to achieve the same result

  smartling-cli files download "**/*.json" --locale fr-FR --locale de-DE

# Download all translated files

  smartling-cli files download --all

# Pull every file × target locale in a translation job

  smartling-cli --threads 20 files pull --job-uid <jobUid>

# Pull only .txt files from a job (URI glob filters the job file list)

  smartling-cli files pull "**.txt" --job-uid <jobUid>

# Preview what a job pull would download

  smartling-cli files pull --job-uid <jobUid> --dry-run
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv(ctx)
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			var threadsCfg *string
			if cnf, cnfErr := rootcmd.Config(); cnfErr == nil && cnf.Threads > 0 {
				s := strconv.FormatUint(uint64(cnf.Threads), 10)
				threadsCfg = new(s)
			}
			threadsParam := resolve.FallbackString(cmd.Flags().Lookup(threadsFlag), resolve.StringParam{
				FlagName: threadsFlag,
				Config:   threadsCfg,
			})
			threadsParamI, err := strconv.ParseUint(threadsParam, 10, 32)
			if err != nil {
				rlog.Errorf("failed to parse `threads` parameter: %s", err)
				os.Exit(1)
			}

			params := files.PullParams{
				URI:       uri,
				JobUID:    jobUID,
				All:       all,
				Format:    formatPath,
				Directory: directory,
				Source:    source,
				Locales:   locales,
				Progress:  progress,
				Retrieve:  retrieve,
				Resume:    resume,
				DryRun:    dryRun,
				Threads:   uint32(threadsParamI),
			}
			err = s.RunPull(ctx, params)
			if err != nil {
				rlog.Errorf("failed to run pull: %s", err)
				os.Exit(1)
			}
		},
	}

	pullCmd.Flags().BoolVar(&all, "all", false, `Download all files. Required if no file pattern is specified.`)
	pullCmd.Flags().StringVar(&jobUID, "job-uid", "", "Filter downloads to files belonging to the specified job UID")
	pullCmd.Flags().BoolVar(&source, "source", false, `Pulls source file as well.`)
	pullCmd.Flags().StringVar(&progress, "progress", "", `Pulls only translations that are at least specified percent of work complete.`)
	pullCmd.Flags().StringVar(&retrieve, "retrieve", "", `Retrieval type: pending, published, pseudo or contextMatchingInstrumented.`)
	pullCmd.Flags().StringVarP(&directory, "directory", "d", ".", `Download all files to specified directory.`)
	pullCmd.Flags().StringArrayVarP(&locales, "locale", "l", []string{}, `Authorize only specified locales.`)
	pullCmd.Flags().BoolVar(&resume, "resume", false, `Resume a previously interrupted pull operation, skipping already downloaded files.`)
	pullCmd.Flags().BoolVar(&dryRun, "dry-run", false, `Print the file × locale matrix that would be downloaded, then exit.`)
	pullCmd.Flags().Uint32Var(&threads, threadsFlag, 20, `If command can be executed concurrently, it will be
executed for at most <number> of threads.`)
	pullCmd.Flags().StringVar(&formatPath, "format", "", `Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           Default: `+format.DefaultFilePullFormat)

	return pullCmd
}
