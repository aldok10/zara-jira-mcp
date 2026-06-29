package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// HealthStatus represents the health status of the server.
type HealthStatus string

const (
	HealthOK       HealthStatus = "ok"
	HealthDegraded HealthStatus = "degraded"
	HealthFailed   HealthStatus = "failed"
)

// HealthInfo contains detailed health information.
type HealthInfo struct {
	Status  HealthStatus   `json:"status"`
	Version string         `json:"version"`
	Uptime  string         `json:"uptime"`
	Modules []ModuleHealth `json:"modules"`
}

// ModuleHealth represents the health of a specific module.
type ModuleHealth struct {
	Name   string       `json:"name"`
	Status HealthStatus `json:"status"`
	Error  string       `json:"error,omitempty"`
}

// HealthChecker provides health checking capabilities.
type HealthChecker struct {
	startTime time.Time
	version   string
	modules   []ModuleHealth
}

// NewHealthChecker creates a new HealthChecker.
func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		startTime: time.Now(),
		version:   version,
		modules:   make([]ModuleHealth, 0),
	}
}

// CheckModule checks the health of a module.
func (h *HealthChecker) CheckModule(ctx context.Context, name string, checkFn func(ctx context.Context) error) ModuleHealth {
	module := ModuleHealth{Name: name, Status: HealthOK}

	if err := checkFn(ctx); err != nil {
		module.Status = HealthFailed
		module.Error = err.Error()
		slog.Error("module health check failed", "module", name, "error", err)
	}

	h.modules = append(h.modules, module)
	return module
}

// GetHealth returns the overall health status.
func (h *HealthChecker) GetHealth(ctx context.Context) HealthInfo {
	info := HealthInfo{
		Status:  HealthOK,
		Version: h.version,
		Uptime:  time.Since(h.startTime).Truncate(time.Second).String(),
		Modules: h.modules,
	}

	// Determine overall status based on modules
	for _, m := range h.modules {
		if m.Status == HealthFailed {
			info.Status = HealthFailed
			break
		}
		if m.Status == HealthDegraded {
			info.Status = HealthDegraded
		}
	}

	return info
}

// GetHealthJSON returns health status as JSON.
func (h *HealthChecker) GetHealthJSON(ctx context.Context) (string, error) {
	info := h.GetHealth(ctx)
	data, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("marshal health: %w", err)
	}
	return string(data), nil
}

// RecordModuleHealth records a module's health status.
func (h *HealthChecker) RecordModuleHealth(name string, status HealthStatus, errMsg string) {
	module := ModuleHealth{
		Name:   name,
		Status: status,
		Error:  errMsg,
	}
	h.modules = append(h.modules, module)
}
