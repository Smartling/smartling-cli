package joblist

import (
	"fmt"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/Smartling/api-sdk-go/helpers/uid"
	"github.com/spf13/cobra"
)

const (
	accountFlag       = "account"
	nameFlag          = "name"
	numberFlag        = "number"
	statusFlag        = "status"
	uidFlag           = "uid"
	projectIDFlag     = "project-id"
	fileFlag          = "file"
	hashcodeFlag      = "hashcode"
	withPriorityFlag  = "with-priority"
	sortByFlag        = "sort-by"
	sortDirectionFlag = "sort-direction"
	limitFlag         = "limit"
	offsetFlag        = "offset"
)

// resolveParams resolves jobs-list params from flags with an env-var
// fallback (flag → env).
func resolveParams(cmd *cobra.Command, projectID, accountIDConfig string) (srv.ListParams, error) {
	account, _ := cmd.Flags().GetBool(accountFlag)

	accountUID := accountIDConfig
	if account {
		resolved, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), accountIDConfig)
		if err != nil {
			return srv.ListParams{}, err
		}
		accountUID = string(resolved)
	}

	// FallbackBool defaults to true when the flag is present-but-unchanged
	// with no env/config, so pass an explicit false default.
	withPriorityDefault := false

	limit, err := resolve.FallbackInt(cmd.Flags().Lookup(limitFlag), resolve.IntParam{FlagName: limitFlag})
	if err != nil {
		return srv.ListParams{}, fmt.Errorf("invalid value for --%s: %w", limitFlag, err)
	}
	offset, err := resolve.FallbackInt(cmd.Flags().Lookup(offsetFlag), resolve.IntParam{FlagName: offsetFlag})
	if err != nil {
		return srv.ListParams{}, fmt.Errorf("invalid value for --%s: %w", offsetFlag, err)
	}

	return srv.ListParams{
		AccountUID:         uid.AccountUID(accountUID),
		ProjectUID:         projectID,
		Account:            account,
		JobName:            resolve.FallbackString(cmd.Flags().Lookup(nameFlag), resolve.StringParam{FlagName: nameFlag}),
		JobNumber:          resolve.FallbackString(cmd.Flags().Lookup(numberFlag), resolve.StringParam{FlagName: numberFlag}),
		JobStatus:          resolve.FallbackStringArray(cmd, statusFlag, nil),
		TranslationJobUIDs: resolve.FallbackStringArray(cmd, uidFlag, nil),
		ProjectIDs:         resolve.FallbackStringArray(cmd, projectIDFlag, nil),
		FileURIs:           resolve.FallbackStringArray(cmd, fileFlag, nil),
		Hashcodes:          resolve.FallbackStringArray(cmd, hashcodeFlag, nil),
		WithPriority:       resolve.FallbackBool(cmd.Flags().Lookup(withPriorityFlag), resolve.BoolParam{FlagName: withPriorityFlag, Config: &withPriorityDefault}),
		SortBy:             resolve.FallbackString(cmd.Flags().Lookup(sortByFlag), resolve.StringParam{FlagName: sortByFlag}),
		SortDirection:      resolve.FallbackString(cmd.Flags().Lookup(sortDirectionFlag), resolve.StringParam{FlagName: sortDirectionFlag}),
		Limit:              uint32(limit),
		Offset:             uint32(offset),
	}, nil
}
