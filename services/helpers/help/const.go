package help

const (
	AuthenticationOptions = `
  --user <user>
    Specify user ID for authentication.

  --secret <secret>
    Specify secret token for authentication.

  -a --account <account>
    Specify account ID.
`

	FormatOption = `
This command supports advanced formatting via --format flag with full
support of Golang templates (https://golang.org/pkg/text/template).

Special formatting functions are available:

  > {{name <variable>}} — return file URI without extension for specified
    <variable>;
  > {{ext <variable}} — return extension from file URI for specified <variable>;
`

	GlobPattern = `argument support globbing with following patterns:

  > ** — matches any number of any chars;
  > *  — matches any number of chars except '/';
  > ?  — matches any single char except '/';
  > [xyz]   — matches 'x', 'y' or 'z' charachers;
  > [!xyz]  — matches not 'x', 'y' or 'z' charachers;
  > {a,b,c} — matches alternatives a, b or c;`
)
