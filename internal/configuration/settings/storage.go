package settings

import (
	"fmt"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Storage contains settings to configure the storage.
type Storage struct {
	// Filepath is the path to the servers.json file. An empty string disables on-disk storage.
	Filepath *string
	// UpdateFilepath is the path to the update-servers.json file.
	UpdateFilepath *string
}

func (s Storage) validate() (err error) {
	if *s.Filepath != "" { // optional
		_, err := filepath.Abs(*s.Filepath)
		if err != nil {
			return fmt.Errorf("filepath is not valid: %w", err)
		}
	}
	if *s.UpdateFilepath != "" { // optional
		_, err := filepath.Abs(*s.UpdateFilepath)
		if err != nil {
			return fmt.Errorf("update filepath is not valid: %w", err)
		}
	}
	return nil
}

func (s *Storage) copy() (copied Storage) {
	return Storage{
		Filepath:       gosettings.CopyPointer(s.Filepath),
		UpdateFilepath: gosettings.CopyPointer(s.UpdateFilepath),
	}
}

func (s *Storage) overrideWith(other Storage) {
	s.Filepath = gosettings.OverrideWithPointer(s.Filepath, other.Filepath)
	s.UpdateFilepath = gosettings.OverrideWithPointer(s.UpdateFilepath, other.UpdateFilepath)
}

func (s *Storage) setDefaults() {
	s.Filepath = gosettings.DefaultPointer(s.Filepath, constants.ServersData)
	s.UpdateFilepath = gosettings.DefaultPointer(s.UpdateFilepath, constants.ServersUpdateData)
}

func (s Storage) String() string {
	return s.toLinesNode().String()
}

func (s Storage) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Storage settings:")
	if *s.Filepath != "" {
		node.Appendf("Filepath: %s", *s.Filepath)
	} else {
		node.Appendf("Filepath: disabled")
	}
	if *s.UpdateFilepath != "" {
		node.Appendf("Update filepath: %s", *s.UpdateFilepath)
	} else {
		node.Appendf("Update filepath: disabled")
	}
	return node
}

func (s *Storage) read(r *reader.Reader) (err error) {
	s.Filepath = r.Get("STORAGE_FILEPATH", reader.AcceptEmpty(true))
	s.UpdateFilepath = r.Get("STORAGE_UPDATE_FILEPATH", reader.AcceptEmpty(true))
	return nil
}
