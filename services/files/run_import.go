package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

type ImportParams struct {
	URI             string
	File            string
	Locale          string
	FileType        string
	PostTranslation bool
	Overwrite       bool
}

func (s Service) RunImport(params ImportParams) error {
	contents, err := os.ReadFile(params.File)
	if err != nil {
		return clierror.NewError(
			hierr.Errorf(err, "unable to read file for import"),
			"Check that specified file exists and you have permissions "+
				"to read it.",
		)
	}

	request := sdk.ImportRequest{}

	request.File = contents
	request.FileURI = params.URI

	request.TranslationState = sdk.TranslationStatePublished

	if params.PostTranslation {
		request.TranslationState = sdk.TranslationStatePostTranslation
	}

	request.Overwrite = params.Overwrite

	if params.FileType != "" {
		request.FileType = sdk.FileType(params.FileType)
	} else {
		request.FileType = sdk.GetFileTypeByExtension(
			filepath.Ext(params.File),
		)

		if request.FileType == sdk.FileTypeUnknown {
			return clierror.NewError(
				fmt.Errorf(
					"unable to deduce file type from extension: %q",
					filepath.Ext(params.File),
				),

				`You need to specify file type via --type option.`,
			)
		}
	}

	result, err := s.Client.Import(s.Config.ProjectID, params.Locale, request)
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
