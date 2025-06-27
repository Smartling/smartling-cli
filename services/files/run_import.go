package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"

	smfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// ImportParams holds the parameters for the RunImport method.
type ImportParams struct {
	URI             string
	File            string
	Locale          string
	FileType        string
	PostTranslation bool
	Overwrite       bool
}

// RunImport imports a file into the Smartling project with the specified parameters.
func (s service) RunImport(params ImportParams) error {
	contents, err := os.ReadFile(params.File)
	if err != nil {
		return clierror.NewError(
			hierr.Errorf(err, "unable to read file for import"),
			"Check that specified file exists and you have permissions "+
				"to read it.",
		)
	}

	request := smfile.ImportRequest{}

	request.File = contents
	request.FileURI = params.URI

	request.TranslationState = smfile.TranslationStatePublished

	if params.PostTranslation {
		request.TranslationState = smfile.TranslationStatePostTranslation
	}

	request.Overwrite = params.Overwrite

	if params.FileType != "" {
		request.FileType = smfile.FileType(params.FileType)
	} else {
		request.FileType = smfile.GetFileTypeByExtension(
			filepath.Ext(params.File),
		)

		if request.FileType == smfile.FileTypeUnknown {
			return clierror.NewError(
				fmt.Errorf(
					"unable to deduce file type from extension: %q",
					filepath.Ext(params.File),
				),

				`You need to specify file type via --type option.`,
			)
		}
	}

	result, err := s.APIClient.Import(s.Config.ProjectID, params.Locale, request)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to import file "%s" (original "%s")`,
			params.File,
			params.URI,
		)
	}

	fmt.Printf(
		"%s imported [%d strings %d words]\n",
		params.File,
		result.StringCount,
		result.WordCount,
	)

	return nil
}
