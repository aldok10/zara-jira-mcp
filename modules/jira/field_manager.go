package jira

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/jira/infrastructure/client"
	"github.com/aldok10/zara-jira-mcp/modules/jira/infrastructure/repository"
)

// FieldManager orchestrates dynamic field selection for Jira issues using board repository and REST client.
type FieldManager struct {
	BoardRepo     repository.BoardRepository
	Rest          client.RestInterface
	DefaultFields []string // Essential fields like "summary", "status", "priority"
}

// NewFieldManager initializes a new dynamic field manager with board repository and REST client.
func NewFieldManager(boardRepo repository.BoardRepository, restClient client.RestInterface) *FieldManager {
	return &FieldManager{
		BoardRepo:     boardRepo,
		Rest:          restClient,
		DefaultFields: []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated", "duedate", "sprint", "story_points"},
	}
}

// GetFieldsForBoard dynamically determines fields to include for a board ID.
func (fm *FieldManager) GetFieldsForBoard(ctx context.Context, boardID int) ([]string, error) {
	if fm == nil {
		return nil, fmt.Errorf("FieldManager is not initialized")
	}

	// Check board repository for cached custom fields
	if cfg, err := fm.BoardRepo.FindByID(ctx, boardID); err == nil && cfg != nil && len(cfg.CustomFields) > 0 {
		return append(fm.DefaultFields, cfg.CustomFields...), nil
	}

	// If no cached fields, attempt to fetch from Jira REST API
	fields, err := fm.Rest.GetCustomFieldsByBoard(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom fields for board %d: %w", boardID, err)
	}

	// Cache the fetched custom fields in the board config
	if cfg, err := fm.BoardRepo.FindByID(ctx, boardID); err == nil && cfg != nil {
		cfg.CustomFields = fields
		fm.BoardRepo.Save(ctx, cfg)
	}

	return append(fm.DefaultFields, fields...), nil
}

// ValidateFieldExists checks if a custom field ID is valid for a specific board.
func (fm *FieldManager) ValidateFieldExists(ctx context.Context, boardID int, fieldID string) (bool, error) {
	if fm == nil {
		return false, fmt.Errorf("FieldManager is not initialized")
	}

	// Check default fields first
	for _, defaultField := range fm.DefaultFields {
		if fieldID == defaultField {
			return true, nil
		}
	}

	// Try to fetch fresh custom fields from repository
	if cfg, err := fm.BoardRepo.FindByID(ctx, boardID); err == nil && cfg != nil {
		for _, customField := range cfg.CustomFields {
			if fieldID == customField {
				return true, nil
			}
		}
	}

	// As a last resort, try to get custom fields from Jira
	fields, err := fm.Rest.GetCustomFieldsByBoard(ctx, boardID)
	if err != nil {
		return false, err
	}

	// Validate against fetched fields
	for _, fetchedField := range fields {
		if fieldID == fetchedField {
			return true, nil
		}
	}

	return false, nil
}

// ManageCustomField synchronizes custom field definitions for a board, either adding or removing them.
func (fm *FieldManager) ManageCustomField(ctx context.Context, boardID int, fieldID string, action string) error {
	if fm == nil {
		return fmt.Errorf("FieldManager is not initialized")
	}

	// Find the board configuration
	cfg, err := fm.BoardRepo.FindByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("unable to find board configuration for boardID %d: %w", boardID, err)
	}

	// Perform action on custom fields
	switch action {
	case "add":
		if !cfg.HasCustomField(fieldID) {
			cfg.AddCustomField(fieldID)
			fm.BoardRepo.Save(ctx, cfg)
		}
	case "remove":
		if cfg.HasCustomField(fieldID) {
			cfg.RemoveCustomField(fieldID)
			fm.BoardRepo.Save(ctx, cfg)
		}
	default:
		return fmt.Errorf("invalid action '%s', expected 'add' or 'remove'", action)
	}

	return nil
}

// GetAllCustomFieldsAcrossBoards retrieves all unique custom fields from all configured boards.
func (fm *FieldManager) GetAllCustomFieldsAcrossBoards(ctx context.Context) ([]string, error) {
	// Since the repository does not expose a way to list all boards, this method requires
	// manual extension or implementation. We'll return an error indicating the limitation.
	return nil, fmt.Errorf("cannot efficiently collect all custom fields across boards without a method to list all boards")
}
