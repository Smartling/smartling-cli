## smartling-cli files push

Creates job and uploads specified file into this job.

### Synopsis

smartling-cli files push <file> [<uri>] --job <job name> [--authorize] [--locale <locale>] [--type <type>] [--branch (@auto|<branch name>)] [--directory <work dir>] [--directive <smartling directive>]

Creates a new job (or reuses existing) in Smartling TMS and uploads designated
file(s) for translation.

One or more files can be pushed.

When pushing single file, <uri> can be specified to override local path.
When pushing multiple files, they will be uploaded using local path as URI.
If no file specified in command line, config file will be used to lookup
for file masks to push.

Use --job option to specify job name or job UID. If job name is not specified,
then the "CLI uploads" name will be used.
You can use the same job name for multiple CLI calls. The same job will be used
in this case (CLI searches by the job name). If the job with the same name exists,
but it has state Canceled or Closed, then a new job will be created with timestamp suffix.

To authorize the job after uploading all files, use --authorize option.

To specify locales for the files in the job, use one or more --locale options.
If no locales specified, then all project locales will be added to all uploaded files.

To prepend prefix to all target URIs, use --branch option. Special
value "@auto" can be used to tell the tool to use the current git
branch name as value for --branch option.

File type will be deduced from file extension. If file extension is unknown,
type should be specified manually by using --type option. That option also
can be used to override detected file type.

`<file>` argument supports globbing with following patterns:

  > ** — matches any number of any chars;
  > *  — matches any number of chars except '/';
  > ?  — matches any single char except '/';
  > [xyz]   — matches 'x', 'y' or 'z' charachers;
  > [!xyz]  — matches not 'x', 'y' or 'z' charachers;
  > {a,b,c} — matches alternatives a, b or c; 


  --user <user>
    Specify user ID for authentication.

  --secret <secret>
    Specify secret token for authentication.

  -a --account <account>
    Specify account ID.


```
smartling-cli files push <file> <uri> --job <job name> [--authorize] [--locale <locale>] [flags]
```

### Examples

```

# Upload files to a translation job

  smartling-cli files push my-file.txt --job "Website Update" --authorize

# Upload multiple files using pattern matching with the command alias ‘upload’

  smartling-cli files upload "src/**/*.json" --job "App Localization"

# Manual branch naming

  smartling-cli files push "**/*.txt" --branch "feature-branch"

# Automatic Git branch detection

  smartling-cli files push "**/*.txt" --branch "@auto"

# All JSON files in subdirectories

  smartling-cli files push "**/*.json"

# Specific file types

  smartling-cli files push "**/*.{json,xml,properties}"

# Files matching naming convention with the command alias 'upload' 

  smartling-cli files upload "**/messages_*.properties"


```

### Options

```
  -z, --authorize               Automatically authorize the job with file(s) and specified locales.
                                If the flag is not specified, the job remains unauthorized.
  -b, --branch string           <branch>
                                Prepend specified prefix to target file URI.
  -r, --directive stringArray   Specify one or more directives to use in push request.
  -d, --directory string        Specified directory. (default ".")
  -h, --help                    help for push
  -j, --job string              <job name>
                                Provide a name for the Smartling translation job or job UID.
                                All files will be uploaded into this job.
                                If the flag is not specified then the "CLI uploads" name will be used.
  -l, --locale stringArray      <locale code>
                                Add file(s) to the job for the specified locale only.
                                If the flag is not specified, then all project locales will be added to the job.
                                Can be specified several times: --locale fr --locale de -l es
  -t, --type string             <type>
                                Override automatically detected file type.
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

###### Auto generated by spf13/cobra on 28-Oct-2025
