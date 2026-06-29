package agent

import (
	"context"
	"time"
)

// Session represents a persistent conversation between a user and the agent.
// Mirrors Hermes session storage concept — per-channel, per-user isolation.
type Session struct {
	ID           string    `json:"id"`
	Channel      string    `json:"channel"`    // lark, slack, discord, telegram
	ChannelID    string    `json:"channel_id"` // chat_id, channel_id
	UserID       string    `json:"user_id"`    // platform user ID
	UserName     string    `json:"user_name"`  // display name
	Summary      string    `json:"summary"`    // running summary for context compression
	MessageCount int       `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Metadata     string    `json:"metadata"` // JSON blob for extra state
}

// Message is a single turn in a conversation.
type Message struct {
	ID         int64     `json:"id"`
	SessionID  string    `json:"session_id"`
	Role       string    `json:"role"`       // user, assistant, tool
	Content    string    `json:"content"`    // text content or tool result
	ToolName   string    `json:"tool_name"`  // if role=tool
	ToolInput  string    `json:"tool_input"` // JSON of tool params
	TokenCount int       `json:"token_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// Memory represents a single remembered fact stored by the user or AI.
type Memory struct {
	ID        int64     `json:"id"`
	SessionID string    `json:"session_id"`
	Scope     string    `json:"scope"` // "session" or "global"
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScheduledTask represents a periodic agent task (cron job).
type ScheduledTask struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Interval  string     `json:"interval"`   // "30m", "1h", "2h", "daily:09:00", "weekday:09:00"
	Prompt    string     `json:"prompt"`     // what to ask the agent
	Channel   string     `json:"channel"`    // lark, slack, etc.
	ChannelID string     `json:"channel_id"` // target channel/chat
	Enabled   bool       `json:"enabled"`
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	NextRunAt time.Time  `json:"next_run_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// SessionStore persists conversations across restarts.
type SessionStore interface {
	// GetOrCreateSession finds an existing session or creates a new one.
	GetOrCreateSession(ctx, channel, channelID, userID, userName string) (*Session, error)

	// GetSession retrieves a session by ID.
	GetSession(ctx, sessionID string) (*Session, error)

	// AppendMessage stores a message and updates session.
	AppendMessage(ctx, sessionID string, msg *Message) error

	// GetRecentMessages returns the last N messages for context.
	GetRecentMessages(ctx, sessionID string, limit int) ([]Message, error)

	// UpdateSessionSummary updates the running conversation summary.
	UpdateSessionSummary(ctx, sessionID, summary string) error

	// SaveMemory stores a fact. Replaces existing if key + session_id exists.
	SaveMemory(ctx, sessionID, key, value, scope string) error

	// GetMemoriesBySession returns all memories for a session.
	GetMemoriesBySession(ctx, sessionID string) ([]Memory, error)

	// GetMemoryByKey returns a specific memory by key.
	GetMemoryByKey(ctx, sessionID, key string) (*Memory, error)

	// DeleteMemory removes a memory by key.
	DeleteMemory(ctx, sessionID, key string) error

	// SearchMemories searches memories by text across key and value.
	SearchMemories(ctx, sessionID, query string) ([]Memory, error)

	// SaveScheduledTask stores a new scheduled task.
	SaveScheduledTask(ctx, task *ScheduledTask) error

	// GetScheduledTasks returns all scheduled tasks, optionally filtered by channel.
	GetScheduledTasks(ctx, channel string) ([]ScheduledTask, error)

	// GetDueTasks returns tasks where next_run_at <= now and enabled = true.
	GetDueTasks(ctx context.Context, now time.Time) ([]ScheduledTask, error)

	// UpdateTaskNextRun sets the next_run_at and last_run_at for a task.
	UpdateTaskRun(ctx, taskID int64, lastRun, nextRun time.Time) error

	// DeleteScheduledTask removes a scheduled task.
	DeleteScheduledTask(ctx, taskID int64) error

	// Close cleanly shuts down the store.
	Close() error
}
