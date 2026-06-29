package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// GetProjects returns all accessible projects.
func (c *RestClient) GetProjects(ctx context.Context) ([]domain.Project, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, "/rest/api/3/project/search?maxResults=50", nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Values []struct {
			Key  string `json:"key"`
			Name string `json:"name"`
			Lead struct {
				DisplayName string `json:"displayName"`
			} `json:"lead"`
			ProjectTypeKey string `json:"projectTypeKey"`
		} `json:"values"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode projects: %w", err)
	}
	out := make([]domain.Project, len(result.Values))
	for i, p := range result.Values {
		out[i] = domain.Project{Key: p.Key, Name: p.Name, Lead: p.Lead.DisplayName, Type: p.ProjectTypeKey}
	}
	return out, nil
}

// GetProject returns project details.
func (c *RestClient) GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/project/%s", key), nil)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Lead struct {
			DisplayName string `json:"displayName"`
		} `json:"lead"`
		ProjectTypeKey string `json:"projectTypeKey"`
		Description    string `json:"description"`
		Components     []struct {
			Name string `json:"name"`
		} `json:"components"`
		Versions []struct {
			Name string `json:"name"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode project: %w", err)
	}
	pd := &domain.ProjectDetail{
		Key:         raw.Key,
		Name:        raw.Name,
		Lead:        raw.Lead.DisplayName,
		Type:        raw.ProjectTypeKey,
		Description: raw.Description,
	}
	for _, comp := range raw.Components {
		pd.Components = append(pd.Components, comp.Name)
	}
	for _, v := range raw.Versions {
		pd.Versions = append(pd.Versions, v.Name)
	}
	return pd, nil
}

// GetFields returns all Jira fields.
func (c *RestClient) GetFields(ctx context.Context) ([]domain.Field, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, "/rest/api/3/field", nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Custom bool   `json:"custom"`
		Schema *struct {
			Type string `json:"type"`
		} `json:"schema"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode fields: %w", err)
	}
	out := make([]domain.Field, len(raw))
	for i, f := range raw {
		ft := ""
		if f.Schema != nil {
			ft = f.Schema.Type
		}
		out[i] = domain.Field{ID: f.ID, Name: f.Name, Custom: f.Custom, Type: ft}
	}
	return out, nil
}

// GetAttachments returns attachments for an issue.
func (c *RestClient) GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error) {
	data, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/rest/api/3/issue/%s?fields=attachment", issueKey), nil)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Fields struct {
			Attachment []struct {
				ID       string `json:"id"`
				Filename string `json:"filename"`
				Size     int64  `json:"size"`
				MimeType string `json:"mimeType"`
				Author   struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Created string `json:"created"`
				Content string `json:"content"`
			} `json:"attachment"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]domain.Attachment, len(raw.Fields.Attachment))
	for i, a := range raw.Fields.Attachment {
		out[i] = domain.Attachment{
			ID: a.ID, Filename: a.Filename, Size: a.Size,
			MimeType: a.MimeType, Author: a.Author.DisplayName,
			Created: a.Created, URL: a.Content,
		}
	}
	return out, nil
}

// GetVersions returns versions for a project.
func (c *RestClient) GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/project/%s/versions", projectKey), nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Released    bool   `json:"released"`
		ReleaseDate string `json:"releaseDate"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode versions: %w", err)
	}
	out := make([]domain.Version, len(raw))
	for i, v := range raw {
		out[i] = domain.Version{ID: v.ID, Name: v.Name, Description: v.Description, Released: v.Released, ReleaseDate: v.ReleaseDate}
	}
	return out, nil
}

// CreateVersion creates a new version.
func (c *RestClient) CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error) {
	payload, err := json.Marshal(map[string]string{"name": name, "description": description, "project": projectKey})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/rest/api/3/version", payload)
	if err != nil {
		return nil, err
	}
	var raw struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode version: %w", err)
	}
	return &domain.Version{ID: raw.ID, Name: raw.Name, Description: description}, nil
}

// ReleaseVersion marks a version as released.
func (c *RestClient) ReleaseVersion(ctx context.Context, versionID string) error {
	payload, err := json.Marshal(map[string]any{"released": true, "releaseDate": "2006-01-02"})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/version/%s", versionID), payload)
	return err
}

// GetComponents returns components for a project.
func (c *RestClient) GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/project/%s/components", projectKey), nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Lead *struct {
			DisplayName string `json:"displayName"`
		} `json:"lead"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode components: %w", err)
	}
	out := make([]domain.Component, len(raw))
	for i, comp := range raw {
		lead := ""
		if comp.Lead != nil {
			lead = comp.Lead.DisplayName
		}
		out[i] = domain.Component{ID: comp.ID, Name: comp.Name, Lead: lead}
	}
	return out, nil
}
