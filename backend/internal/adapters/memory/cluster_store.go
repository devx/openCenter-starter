package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/devx/openCenter-starter/backend/internal/ports"
)

type ClusterStore struct {
	mu       sync.RWMutex
	clusters map[string]ports.Cluster
}

func NewClusterStore() *ClusterStore {
	return &ClusterStore{
		clusters: map[string]ports.Cluster{},
	}
}

func (s *ClusterStore) List(_ context.Context, filter ports.ClusterFilter) ([]ports.Cluster, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ports.Cluster, 0, len(s.clusters))
	for _, cluster := range s.clusters {
		if filter.Status != "" && cluster.Status != filter.Status {
			continue
		}
		if filter.NamePrefix != "" && !strings.HasPrefix(cluster.Name, filter.NamePrefix) {
			continue
		}
		if filter.IDPrefix != "" && !strings.HasPrefix(cluster.ID, filter.IDPrefix) {
			continue
		}
		result = append(result, cluster)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	total := len(result)
	if filter.Offset >= total {
		return []ports.Cluster{}, total, nil
	}

	end := filter.Offset + filter.Limit
	if end > total {
		end = total
	}

	return result[filter.Offset:end], total, nil
}

func (s *ClusterStore) Get(_ context.Context, id string) (ports.Cluster, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cluster, ok := s.clusters[id]
	return cluster, ok, nil
}

func (s *ClusterStore) Create(_ context.Context, name string) (ports.Cluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cluster := ports.Cluster{
		ID:     uuid.NewString(),
		Name:   name,
		Status: "provisioning",
	}
	s.clusters[cluster.ID] = cluster
	return cluster, nil
}

func (s *ClusterStore) Update(_ context.Context, id string, update ports.ClusterUpdate) (ports.Cluster, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cluster, ok := s.clusters[id]
	if !ok {
		return ports.Cluster{}, false, nil
	}

	if update.Name != nil {
		cluster.Name = *update.Name
	}
	if update.Status != nil {
		cluster.Status = *update.Status
	}

	s.clusters[id] = cluster
	return cluster, true, nil
}

func (s *ClusterStore) Archive(_ context.Context, id string) (ports.Cluster, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cluster, ok := s.clusters[id]
	if !ok {
		return ports.Cluster{}, false, nil
	}

	cluster.Status = "archived"
	s.clusters[id] = cluster
	return cluster, true, nil
}
