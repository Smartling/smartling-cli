package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/99designs/smartling"
	"github.com/99designs/smartling/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var ProjectCommand = cli.Command{
	Name:  "project",
	Usage: "manage local project files",
	Before: func(c *cli.Context) error {
		err := cmdBefore(c)
		if err != nil {
			return nil
		}

		return loadProjectErr
	},
	Subcommands: []cli.Command{
		projectFilesCommand,
		projectStatusCommand,
		projectPullCommand,
		projectPushCommand,
	},
}

func fetchRemoteFileList() stringSlice {
	files := stringSlice{}
	listFiles, err := client.List(smartling.ListRequest{})
	logAndQuitIfError(err)

	for _, fs := range listFiles {
		files = append(files, fs.FileUri)
	}

	return files
}

func fetchLocales() []string {
	ll := []string{}
	locales, err := client.Locales()
	logAndQuitIfError(err)
	for _, l := range locales {
		ll = append(ll, l.Locale)
	}

	return ll
}

var projectFilesCommand = cli.Command{
	Name:        "files",
	Usage:       "lists the local files",
	Description: "files",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 0 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: files")
		}

		for _, projectFilepath := range ProjectConfig.Files() {
			fmt.Println(projectFilepath)
		}
	},
}

var projectStatusCommand = cli.Command{
	Name:        "status",
	Usage:       "show the status of the project's remote files",
	Description: "status [<prefix>]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "awaiting-auth",
			Usage: "Output the number of strings Awaiting Authorization",
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: status [<prefix>]")
		}

		prefix := prefixOrGitPrefix(c.Args().Get(0))
		locales := fetchLocales()
		statuses := GetProjectStatus(prefix, locales)

		if c.Bool("awaiting-auth") {
			fmt.Println(statuses.AwaitingAuthorizationCount())
		} else {
			fmt.Print("\n")
			PrintProjectStatusTable(statuses, locales)
			fmt.Print("\n")
			fmt.Printf("Awaiting Authorization: %4d\n", statuses.AwaitingAuthorizationCount())
			fmt.Printf("Total:                  %4d\n", statuses.TotalStringsCount())
		}

	},
}

var projectPullCommand = cli.Command{
	Name:        "pull",
	Usage:       "translate local project files using Smartling as a translation memory",
	Description: "pull [<prefix>]",
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: pull [<prefix>]")
		}

		prefix := prefixOrGitPrefix(c.Args().Get(0))

		pullAllProjectFiles(prefix)
	},
}

func pullAllProjectFiles(prefix string) {
	locales, err := client.Locales()
	logAndQuitIfError(err)

	var wg sync.WaitGroup
	for _, projectFilepath := range ProjectConfig.Files() {
		for _, l := range locales {
			wg.Add(1)
			go func(locale, projectFilepath string) {
				defer wg.Done()

				pullProjectFile(projectFilepath, locale, prefix)
			}(l.Locale, projectFilepath)
		}
	}
	wg.Wait()
}

func pullProjectFile(projectFilepath, locale, prefix string) {
	hit, b, err := translateProjectFile(projectFilepath, locale, prefix)
	logAndQuitIfError(err)

	fp := localPullFilePath(projectFilepath, locale)
	cached := ""
	if hit {
		cached = "(from cache)"
	}
	fmt.Println("Pulled", fp, cached)
	err = ioutil.WriteFile(fp, b, 0644)
	logAndQuitIfError(err)
}

func cleanPrefix(s string) string {
	s = path.Clean("/" + s)
	if s == "/" {
		return ""
	}
	return s
}

func prefixOrGitPrefix(prefix string) string {
	if prefix == "" {
		prefix = pushPrefix()
	}

	prefix = cleanPrefix(prefix)

	if prefix != "" {
		log.Println("Using prefix", prefix)
	}
	return prefix
}

type RemoteFileStatus struct {
	RemoteFilePath string
	Statuses       map[string]*smartling.FileStatus
}

func (r *RemoteFileStatus) NotCompletedStringCount() int {
	c := 0
	for _, fs := range r.Statuses {
		c += fs.NotCompletedStringCount()
	}
	return c
}

func fetchStatusForLocales(remoteFilePath string, locales []string) RemoteFileStatus {
	ss := RemoteFileStatus{
		RemoteFilePath: remoteFilePath,
		Statuses:       map[string]*smartling.FileStatus{},
	}

	var wg sync.WaitGroup
	for _, locale := range locales {
		wg.Add(1)
		go func(f, l string) {
			defer wg.Done()

			s, err := client.Status(f, l)
			logAndQuitIfError(err)
			ss.Statuses[l] = &s

		}(remoteFilePath, locale)
	}
	wg.Wait()

	return ss
}

var projectPushCommand = cli.Command{
	Name:  "push",
	Usage: "upload local project files that contain untranslated strings",
	Description: `push [<prefix>]
Outputs the uploaded files for the given prefix
`,
	Action: func(c *cli.Context) {
		if len(c.Args()) > 1 {
			log.Println("Wrong number of arguments")
			log.Fatalln("Usage: push [<prefix>]")
		}

		prefix := prefixOrGitPrefix(c.Args().Get(0))

		pushAllProjectFiles(prefix)
	},
}

// if prefix is empty, don't append the hash also
func projectFileRemoteName(projectFilepath, prefix string) string {
	remoteFile := projectFilepath
	if prefix != "" {
		remoteFile = fmt.Sprintf("%s/%s/%s", prefix, projectFileHash(projectFilepath), projectFilepath)
	}

	return path.Clean(remoteFile)
}

func pushProjectFile(projectFilepath, prefix string) string {
	remoteFile := projectFileRemoteName(projectFilepath, prefix)

	_, err := client.Upload(projectFilepath, &smartling.UploadRequest{
		FileUri:      remoteFile,
		FileType:     filetypeForProjectFile(projectFilepath),
		ParserConfig: ProjectConfig.ParserConfig,
	})
	logAndQuitIfError(err)

	fmt.Println("Pushed", remoteFile)
	return remoteFile
}

func pushAllProjectFiles(prefix string) {
	var wg sync.WaitGroup
	for _, projectFilepath := range ProjectConfig.Files() {
		wg.Add(1)
		go func(projectFilepath string) {
			defer wg.Done()
			_ = pushProjectFile(projectFilepath, prefix)
		}(projectFilepath)
	}
	wg.Wait()
}

func filetypeForProjectFile(projectFilepath string) smartling.FileType {
	ft := smartling.FileTypeByExtension(path.Ext(projectFilepath))
	if ft == "" {
		ft = ProjectConfig.FileType
	}
	if ft == "" {
		log.Panicln("Can't determine file type for " + projectFilepath)
	}

	return ft
}

type FilenameParts struct {
	Path           string
	Base           string
	Dir            string
	Ext            string
	PathWithoutExt string
	Locale         string
}

func localRelativeFilePath(remotepath string) string {
	fp, err := filepath.Rel(".", path.Join(ProjectConfig.path, remotepath))
	logAndQuitIfError(err)
	return fp
}

func localPullFilePath(p, locale string) string {
	parts := FilenameParts{
		Path:   p,
		Dir:    path.Dir(p),
		Base:   path.Base(p),
		Ext:    path.Ext(p),
		Locale: locale,
	}

	dt := defaultPullDestination
	if dt != "" {
		dt = ProjectConfig.PullFilePath
	}

	out := bytes.NewBufferString("")
	tmpl := template.New("name")
	tmpl.Funcs(template.FuncMap{
		"TrimSuffix": strings.TrimSuffix,
		"Truncate": func(s string, n int) string {
			return s[:n]
		},
	})
	_, err := tmpl.Parse(dt)
	logAndQuitIfError(err)

	err = tmpl.Execute(out, parts)
	logAndQuitIfError(err)

	return localRelativeFilePath(out.String())
}
