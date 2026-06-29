package domain

import "strings"

// Contains adds a custom field to the board configuration if not already present.
func (bc *BoardConfiguration) AddCustomField(fieldID string) {
	for _, existing := range bc.CustomFields {
		if strings.EqualFold(existing, fieldID) {
			return // Already present, no-op
		}
	}
	bc.CustomFields = append(bc.CustomFields, fieldID)
}

// RemoveCustomField removes a field from the board configuration.
func (bc *BoardConfiguration) RemoveCustomField(fieldID string) {
	for i, existing := range bc.CustomFields {
		if strings.EqualFold(existing, fieldID) {
			// Remove using slice tricks
			bc.CustomFields = append(bc.CustomFields[:i], bc.CustomFields[i+1:]...)
			return
		}
	}
}

// GetStatusMapping returns the custom status mapping from the configuration.
// This provides board-specific custom status categorization beyond standard columns.
func (bc *BoardConfiguration) GetStatusMapping(statusName string) (string, bool) {
	// Check for custom status mappings (if added later)
	for status, category := range bc.CustomStatusMappings {
		if strings.EqualFold(status, statusName) {
			return category, true
		}
	}
	return "", false
}

// AllCustomFields returns a copy of the custom fields slice.
func (bc *BoardConfiguration) AllCustomFields() []string {
	if bc == nil {
		return nil
	}
	fields := make([]string, len(bc.CustomFields))
	copy(fields, bc.CustomFields)
	return fields
}

// HasCustomField checks if a field is in the board's custom fields.
func (bc *BoardConfiguration) HasCustomField(fieldID string) bool {
	for _, existing := range bc.CustomFields {
		if strings.EqualFold(existing, fieldID) {
			return true
		}
	}
	return false
}
