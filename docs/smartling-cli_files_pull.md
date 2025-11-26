## smartling-cli files pull

Pulls specified files from server.

### Synopsis

smartling-cli files pull — downloads translated files from project.

Downloads files from specified project into local directory.

It's possible to download only specific files by file mask, to download source
files with translations, to download file to specific directory or to download
specific locales only.

If special value of "-" is specified as <uri>, then program will expect
to read files list from stdin:

  cat files.txt | smartling-cli files pull -

`<uri>` argument supports globbing with following patterns:

  > ** — matches any number of any chars;
  > *  — matches any number of chars except '/';
  > ?  — matches any single char except '/';
  > [xyz]   — matches 'x', 'y' or 'z' charachers;
  > [!xyz]  — matches not 'x', 'y' or 'z' charachers;
  > {a,b,c} — matches alternatives a, b or c; 

If --locale flag is not specified, all available locales are downloaded. To
see available locales, use "status" command.

To download files into subdirectory, use --directory option and specify
directory name you want to download into.

To download source file as well as translated files specify --source option.

Files will be downloaded and stored under names used while upload (e.g. File
URI). While downloading translated file suffix "_<locale>" will be appended to
file name before extension. To override file format name, use --format option.

This command supports advanced formatting via --format flag with full
support of Golang templates (https://golang.org/pkg/text/template).

Special formatting functions are available:

  > {{name <variable>}} — return file URI without extension for specified
    <variable>;
  > {{ext <variable}} — return extension from file URI for specified <variable>;

Following variables are available:

  > .FileURI — full file URI in Smartling system;
  > .Locale — locale ID for translated file and empty for source file;


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

  --user <user>
    Specify user ID for authentication.

  --secret <secret>
    Specify secret token for authentication.

  -a --account <account>
    Specify account ID.


```
smartling-cli files pull <uri> [flags]
```

### Examples

```

# Pull translated files

  smartling-cli files pull "**/*.json" --locale fr-FR --locale de-DE

# Use the alias 'download' to achieve the same result

  smartling-cli files download "**/*.json" --locale fr-FR --locale de-DE

# Download all translated files

  smartling-cli files download --all

```

### Options

```
      --all                  Download all files. Required if no file pattern is specified.
  -d, --directory string     Download all files to specified directory. (default ".")
      --format string        Can be used to format path to downloaded files.
                                                        Note, that single file can be translated in
                                                        different locales, so format should include locale
                                                        to create several file paths.
                                                        Default: {{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}
  -h, --help                 help for pull
  -l, --locale stringArray   Authorize only specified locales.
      --progress string      Pulls only translations that are at least specified percent of work complete.
      --retrieve string      Retrieval type: pending, published, pseudo or contextMatchingInstrumented.
      --source               Pulls source file as well.
```

### Options inherited from parent commands

```
  -a, --account string               Account ID to operate on.
                                     This option overrides config value "account_id".
  -c, --config string                Config file in YAML format.
                                     By default CLI will look for file named
                                     "smartling.yml" in current directory and in all
                                     intermediate parents, emulating git behavior.
  -k, --insecure                     Skip HTTPS certificate validation.
      --operation-directory string   Sets directory to operate on, usually, to store or to
                                     read files.  Depends on command. (default ".")
  -p, --project string               Project ID to operate on.
                                     This option overrides config value "project_id".
      --proxy string                 Use specified URL as proxy server.
      --secret string                Token Secret which will be used for authentication.
                                     This option overrides config value "secret".
      --smartling-url string         Specify base Smartling URL, merely for testing
                                     purposes.
      --threads uint32               If command can be executed concurrently, it will be
                                     executed for at most <number> of threads. (default 4)
      --user string                  User ID which will be used for authentication.
                                     This option overrides config value "user_id".
  -v, --verbose count                Verbose logging
```

### SEE ALSO

* [smartling-cli files](smartling-cli_files.md)	 - Used to access various files sub-commands.

###### Auto generated by spf13/cobra on 26-Nov-2025
