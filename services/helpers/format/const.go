package format

const (
	// DefaultProjectsLocalesFormat is the default format for project locales.
	DefaultProjectsLocalesFormat = `{{.LocaleID}}\t{{.Description}}\t{{.Enabled}}\n`
	// DefaultFilesListFormat is the default format for list files.
	DefaultFilesListFormat = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	// DefaultFileStatusFormat is the default format for file status.
	DefaultFileStatusFormat = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
	// DefaultFilePullFormat is the default format for file pull.
	DefaultFilePullFormat = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
)
