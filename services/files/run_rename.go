package files

import (
	"context"

	"github.com/reconquest/hierr-go"
)

// RunRename renames a file from oldURI to newURI.
func (s service) RunRename(ctx context.Context, oldURI, newURI string) error {
	err := s.APIClient.RenameFile(ctx, s.Config.ProjectID, oldURI, newURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to rename file "%s" -> "%s"`,
			oldURI,
			newURI,
		)
	}
	return nil
}
