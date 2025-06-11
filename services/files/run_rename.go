package files

import (
	"github.com/reconquest/hierr-go"
)

// RunRename renames a file from oldURI to newURI.
func (s Service) RunRename(oldURI, newURI string) error {
	err := s.Client.RenameFile(s.Config.ProjectID, oldURI, newURI)
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
