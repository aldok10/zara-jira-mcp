package repository

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	"github.com/aldok10/zara-jira-mcp/modules/jira/infrastructure/client"
)

// BoardRepository provides board configuration persistence.
type BoardRepository interface {
	Save(ctx context.Context, cfg *domain.BoardConfiguration) error
	FindByID(ctx context.Context, boardID int) (*domain.BoardConfiguration, error)
	Delete(ctx context.Context, boardID int) error
	GetCustomFieldsFromJira(ctx context.Context, boardID int) ([]string, error)
}

// InMemoryBoardRepository is a simple in-memory implementation for prototyping.
type InMemoryBoardRepository struct {
	boards map[int]*domain.BoardConfiguration
	client client.RestInterface
}

func NewInMemoryBoardRepository() *InMemoryBoardRepository {
	return &InMemoryBoardRepository{boards: make(map[int]*domain.BoardConfiguration)}
}

func (r *InMemoryBoardRepository) Save(ctx context.Context, cfg *domain.BoardConfiguration) error {
	r.boards[cfg.ID] = cfg
	return nil
}

func (r *InMemoryBoardRepository) FindByID(ctx context.Context, boardID int) (*domain.BoardConfiguration, error) {
	if cfg, ok := r.boards[boardID]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("board configuration not found")
}

func (r *InMemoryBoardRepository) Delete(ctx context.Context, boardID int) error {
	delete(r.boards, boardID)
	return nil
}

func (r *InMemoryBoardRepository) GetCustomFieldsFromJira(ctx context.Context, boardID int) ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("client not configured")
	}

	// Check cache first
	if cfg, exists := r.boards[boardID]; exists && len(cfg.CustomFields) > 0 {
		return cfg.CustomFields, nil
	}

	// Fetch from Jira and cache
	fields, err := r.client.GetCustomFieldsByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	// Ensure we always have story_points as a standard field
	if !contains(fields, "story_points") {
		fields = append([]string{"story_points"}, fields...)
	}

	// Cache the fields in board config
	if cfg, exists := r.boards[boardID]; exists {
		cfg.CustomFields = fields
	}

	return fields, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
