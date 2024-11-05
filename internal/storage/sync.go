package storage

import (
	"fmt"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
)

func countServers(allServers models.AllServers) (count int) {
	for _, servers := range allServers.ProviderToServers {
		count += len(servers.Servers)
	}
	return count
}

// syncServers merges the hardcoded servers with the ones from the files.
func (s *Storage) syncServers() (err error) {
	hardcodedVersions := make(map[string]uint16, len(s.hardcodedServers.ProviderToServers))
	for provider, servers := range s.hardcodedServers.ProviderToServers {
		hardcodedVersions[provider] = servers.Version
	}

	serversOnFile, err := s.readFromFile(s.filepath, hardcodedVersions)
	if err != nil {
		return fmt.Errorf("reading servers from file: %w", err)
	}

	hardcodedCount := countServers(s.hardcodedServers)
	countOnFile := countServers(serversOnFile)

	s.mergedMutex.Lock()
	defer s.mergedMutex.Unlock()

	if countOnFile == 0 {
		s.logger.Info(fmt.Sprintf(
			"initializing storage %s with %d hardcoded servers",
			s.filepath, hardcodedCount))
		s.mergedServers = s.hardcodedServers
	} else {
		s.logger.Info(fmt.Sprintf(
			"merging by most recent %d hardcoded servers and %d servers read from %s",
			hardcodedCount, countOnFile, s.filepath))
		s.mergedServers = s.mergeServers(s.hardcodedServers, serversOnFile)
	}

	// If the on-disk servers are stale, update the file. This will do nothing if filepath is empty.
	if !reflect.DeepEqual(serversOnFile, s.mergedServers) {
		err = s.flushToFile(s.filepath, nil)
		if err != nil {
			return fmt.Errorf("writing servers to file: %w", err)
		}
	}

	// Update the merged servers with the servers update file.
	serversUpdateFile, err := s.readFromFile(s.updateFilepath, hardcodedVersions)
	if err != nil {
		return fmt.Errorf("failed to read update servers file: %w", err)
	}
	if countServers(serversUpdateFile) > 0 {
		s.mergedServers = s.mergeServers(s.mergedServers, serversUpdateFile)
	}

	return nil
}
