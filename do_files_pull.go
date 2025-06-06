package main

import (
	"github.com/Smartling/api-sdk-go"
	files2 "github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"
	"github.com/Smartling/smartling-cli/services/helpers/thread_pool"
)

func doFilesPull(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		project = config.ProjectID
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = format.DefaultFilePullFormat
	}

	var (
		err   error
		files []smartling.File
	)

	if uri == "-" {
		files, err = reader.ReadFilesFromStdin()
		if err != nil {
			return err
		}
	} else {
		files, err = globfiles.Remote(client, project, uri)
		if err != nil {
			return err
		}
	}

	pool := threadpool.NewThreadPool(config.Threads)

	for _, file := range files {
		// func closure required to pass different file objects to goroutines
		func(file smartling.File) {
			pool.Do(func() {
				err := files2.downloadFileTranslations(client, config, args, file)

				if err != nil {
					logger.Error(err)
				}
			})
		}(file)
	}

	pool.Wait()

	return nil
}
