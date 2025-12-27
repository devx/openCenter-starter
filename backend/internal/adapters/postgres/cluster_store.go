package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devx/openCenter-starter/backend/internal/ports"
)

type ClusterStore struct {
	pool *pgxpool.Pool
}

func NewClusterStore(ctx context.Context, databaseURL string) (*ClusterStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &ClusterStore{pool: pool}, nil
}

func (s *ClusterStore) List(ctx context.Context, filter ports.ClusterFilter) ([]ports.Cluster, int, error) {
	const baseQuery = `
		SELECT id, name, status, created_at, updated_at
		FROM clusters
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR name LIKE $2 || '%')
		  AND ($3 = '' OR id LIKE $3 || '%')
		ORDER BY name
		LIMIT $4 OFFSET $5`

	rows, err := s.pool.Query(ctx, baseQuery, filter.Status, filter.NamePrefix, filter.IDPrefix, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	clusters := make([]ports.Cluster, 0)
	for rows.Next() {
		var cluster ports.Cluster
		if err := rows.Scan(&cluster.ID, &cluster.Name, &cluster.Status, &cluster.CreatedAt, &cluster.UpdatedAt); err != nil {
			return nil, 0, err
		}
		clusters = append(clusters, cluster)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	const countQuery = `
		SELECT COUNT(*)
		FROM clusters
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR name LIKE $2 || '%')
		  AND ($3 = '' OR id LIKE $3 || '%')`

	var total int
	if err := s.pool.QueryRow(ctx, countQuery, filter.Status, filter.NamePrefix, filter.IDPrefix).Scan(&total); err != nil {
		return nil, 0, err
	}

	return clusters, total, nil
}

func (s *ClusterStore) Get(ctx context.Context, id string) (ports.Cluster, bool, error) {
	const query = `
		SELECT id, name, status, created_at, updated_at
		FROM clusters
		WHERE id = $1`

	var cluster ports.Cluster
	err := s.pool.QueryRow(ctx, query, id).Scan(&cluster.ID, &cluster.Name, &cluster.Status, &cluster.CreatedAt, &cluster.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ports.Cluster{}, false, nil
	}
	if err != nil {
		return ports.Cluster{}, false, err
	}

	return cluster, true, nil
}

func (s *ClusterStore) Create(ctx context.Context, name string) (ports.Cluster, error) {
	const query = `
		INSERT INTO clusters (id, name, status)
		VALUES ($1, $2, $3)
		RETURNING id, name, status, created_at, updated_at`

	clusterID := uuid.NewString()
	status := "provisioning"

	var cluster ports.Cluster
	if err := s.pool.QueryRow(ctx, query, clusterID, name, status).Scan(&cluster.ID, &cluster.Name, &cluster.Status, &cluster.CreatedAt, &cluster.UpdatedAt); err != nil {
		return ports.Cluster{}, err
	}

	return cluster, nil
}

func (s *ClusterStore) Update(ctx context.Context, id string, update ports.ClusterUpdate) (ports.Cluster, bool, error) {
	const query = `
		UPDATE clusters
		SET name = COALESCE($2, name),
			status = COALESCE($3, status),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, status, created_at, updated_at`

	var cluster ports.Cluster
	err := s.pool.QueryRow(ctx, query, id, update.Name, update.Status).Scan(&cluster.ID, &cluster.Name, &cluster.Status, &cluster.CreatedAt, &cluster.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ports.Cluster{}, false, nil
	}
	if err != nil {
		return ports.Cluster{}, false, err
	}

	return cluster, true, nil
}

func (s *ClusterStore) Archive(ctx context.Context, id string) (ports.Cluster, bool, error) {
	const query = `
		UPDATE clusters
		SET status = 'archived',
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, status, created_at, updated_at`

	var cluster ports.Cluster
	err := s.pool.QueryRow(ctx, query, id).Scan(&cluster.ID, &cluster.Name, &cluster.Status, &cluster.CreatedAt, &cluster.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ports.Cluster{}, false, nil
	}
	if err != nil {
		return ports.Cluster{}, false, err
	}

	return cluster, true, nil
}
