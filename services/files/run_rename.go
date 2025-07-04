package files

import (
	"github.com/reconquest/hierr-go"
)

// RunRename renames a file from oldURI to newURI.
func (s service) RunRename(oldURI, newURI string) error {
	err := s.APIClient.RenameFile(s.Config.ProjectID, oldURI, newURI)
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
