package ports

import (
	"context"
	"time"
)

type Cluster struct {
	ID        string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClusterUpdate struct {
	Name   *string
	Status *string
}

type ClusterFilter struct {
	Status     string
	NamePrefix string
	IDPrefix   string
	Limit      int
	Offset     int
}

type ClusterStore interface {
	List(ctx context.Context, filter ClusterFilter) ([]Cluster, int, error)
	Get(ctx context.Context, id string) (Cluster, bool, error)
	Create(ctx context.Context, name string) (Cluster, error)
	Update(ctx context.Context, id string, update ClusterUpdate) (Cluster, bool, error)
	Archive(ctx context.Context, id string) (Cluster, bool, error)
}
