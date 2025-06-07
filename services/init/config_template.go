package init

import (
	"strings"
	"text/template"
)

var (
	configTemplate = template.Must(
		template.New(`config`).Delims(`{%`, `%}`).Parse(
			strings.TrimSpace(`
# Config file is optional and all configuration options can be set from command
# line interface.

# (required) Smartling API V2.0 User Identifier used for authentication.
#
# Must be set either in config file or be passed via command line arguments.
user_id:
    "{% .UserID %}"


# (required) Smartling API V2.0 Token Secret used for authentication.
#
# Must be set either in config file or be passed via command line arguments.
secret:
    "{% .Secret %}"


# (optional) Account ID used for projects list requests.
account_id:
    "{% .AccountID %}"


# (optional) Project ID you are working on.
project_id:
    "{% .ProjectID %}"


# (optional) Use specified URL as proxy in case Smartling is not directly
# accessible. An example:
# http://192.168.0.1:3128 - proxy without authentication
# http://user:password@192.168.0.1:3128 - proxy with authentication
#proxy:
#    "PROXY_URL"

# (optional) Additional file-specific settings for push and pull commands.
files:
    # (optional) Special default section will apply configuration to all file
    # types except files, which URIs match following rules.
    default:
        # (optional) Defines pull-specific options.
        pull:
            # (optional) Format, which will be used to format file name.
            #
            # If not set, then default format will be used or format,
            # that is set via command line options.
            format: "{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}"


    # (optional) Specific file settings which uses same pattern rules as CLI
    # tool:
    # > *  - matches everything except /.
    # > ** - matches everything.
    #
    # Note, that pattern should start either with /, * or ** to be matched
    # in case when file was pushed with leading /. Checkout files list in
    # your project first.
    "/path/to/*.properties":
        # (optional) Defines push-specific options.
        push:
            # (optional) Overrides automatically detected file type.
            type: "java_properties"

            # (optional) Sets specific API directives, which are used only
            # for push command. Refer to Smartling API documentation for
            # list of that directives.
            directives:
                namespace: "java"
                file_charset: "utf-8"
                # Any number of custom smartling directives can be specified
                # there.

        pull:
            format: "{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}"

# vim: ft=yaml
`)))
)
