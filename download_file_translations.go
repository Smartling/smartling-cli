package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func downloadFileTranslations(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
	file smartling.File,
) error {
	var (
		project   = config.ProjectID
		directory = args["--directory"].(string)
		source    = args["--source"].(bool)
		locales   = args["--locale"].([]string)

		format, formatGiven = args["--format"].(string)
		progress, _         = args["--progress"].(string)
		retrieve, _         = args["--retrieve"].(string)
	)

	progress = strings.TrimSuffix(progress, "%")
	if progress == "" {
		progress = "0"
	}

	percents, err := strconv.ParseInt(progress, 10, 0)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to parse --progress as integer",
		)
	}

	retrievalType := smartling.RetrievalType(retrieve)

	if format == "" {
		format = defaultFileStatusFormat
	}

	status, err := client.GetFileStatus(project, file.FileURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to retrieve file "%s" locales from project "%s"`,
			file.FileURI,
			project,
		)
	}

	var translations []smartling.FileStatusTranslation

	if source {
		translations = []smartling.FileStatusTranslation{
			{LocaleID: ""},
		}
	} else {
		translations = status.Items
	}

	for _, locale := range translations {
		var complete int64

		if locale.CompletedStringCount > 0 {
			complete = int64(
				100 *
					float64(locale.CompletedStringCount) /
					float64(status.TotalStringCount),
			)
		}

		if percents > 0 {
			if complete < percents {
				continue
			}
		}

		if len(locales) > 0 {
			if !hasLocaleInList(locale.LocaleID, locales) {
				continue
			}
		}

		useFormat := usePullFormat
		if formatGiven {
			useFormat = func(FileConfig) string {
				return format
			}
		}

		path, err := executeFileFormat(
			config,
			file,
			format,
			useFormat,
			map[string]interface{}{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			},
		)
		if err != nil {
			return err
		}

		path = filepath.Join(directory, path)

		err = downloadFile(
			client,
			project,
			file,
			locale.LocaleID,
			path,
			retrievalType,
		)
		if err != nil {
			return err
		}

		if source {
			fmt.Printf("downloaded %s\n", path)
		} else {
			fmt.Printf("downloaded %s %d%%\n", path, int(complete))
		}
	}

	return err
}

func hasLocaleInList(locale string, locales []string) bool {
	for _, filter := range locales {
		if strings.EqualFold(strings.ToLower(filter), strings.ToLower(locale)) {
			return true
		}
	}

	return false
}
