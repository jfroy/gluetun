package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

// FlushToFile flushes the merged servers data to the file
// specified by path, as indented JSON.
func (s *Storage) FlushToFile(path string, providers []string) error {
	if path == "" {
		return nil
	}
	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()
	return s.flushToFile(path, providers)
}

// flushToFile flushes the merged servers data to the file
// specified by path, as indented JSON. It is not thread-safe.
func (s *Storage) flushToFile(path string, providers []string) error {
	if path == "" {
		return nil
	}

	const permission = 0o644
	dirPath := filepath.Dir(path)
	if err := os.MkdirAll(dirPath, permission); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permission)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	servers := s.mergedServers
	if providers != nil {
		servers = models.AllServers{
			Version:           s.mergedServers.Version,
			ProviderToServers: map[string]models.Servers{},
		}
		for _, k := range providers {
			if p, ok := s.mergedServers.ProviderToServers[k]; ok {
				servers.ProviderToServers[k] = p
			}
		}
	}

	for _, obj := range servers.ProviderToServers {
		sort.Sort(models.SortableServers(obj.Servers))
	}

	err = encoder.Encode(&servers)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
