package format

const (
	DefaultProjectsLocalesFormat = `{{.LocaleID}}\t{{.Description}}\t{{.Enabled}}\n`
	DefaultFilesListFormat       = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	DefaultFileStatusFormat      = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
	DefaultFilePullFormat        = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
)
