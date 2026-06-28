package lark

import (
	"context"
	"fmt"
	"strings"
)

func (b *BotHandler) handleCommand(ctx context.Context, text string) string {
	cmd := strings.ToLower(strings.Fields(text)[0])
	switch cmd {
	case "/help", "/h":
		return "/help - Show this message\n/blockers - Active blockers\n/risks - Open risks\n/actions - Pending actions\n/status - Quick summary"
	case "/blockers", "/b":
		return b.cmdBlockers(ctx)
	case "/risks", "/r":
		return b.cmdRisks(ctx)
	case "/actions", "/a":
		return b.cmdActions(ctx)
	case "/status", "/s":
		return b.cmdStatus(ctx)
	default:
		return fmt.Sprintf("Unknown command: %s\nTry /help", cmd)
	}
}

func (b *BotHandler) cmdBlockers(ctx context.Context) string {
	blockers, err := b.memory.GetActiveBlockers(ctx)
	if err != nil {
		return "Failed to fetch blockers."
	}
	if len(blockers) == 0 {
		return "No active blockers."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Active Blockers (%d)\n", len(blockers)))
	for i, bl := range blockers {
		if i >= 10 {
			sb.WriteString(fmt.Sprintf("... and %d more", len(blockers)-10))
			break
		}
		key := bl.IssueKey
		if key == "" {
			key = "-"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", key, bl.Description))
	}
	return sb.String()
}

func (b *BotHandler) cmdRisks(ctx context.Context) string {
	risks, err := b.memory.GetOpenRisks(ctx)
	if err != nil {
		return "Failed to fetch risks."
	}
	if len(risks) == 0 {
		return "No open risks."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Risks (%d)\n", len(risks)))
	for i := range risks {
		if i >= 10 {
			sb.WriteString(fmt.Sprintf("... and %d more", len(risks)-10))
			break
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", strings.ToUpper(risks[i].Severity), risks[i].Title))
	}
	return sb.String()
}

func (b *BotHandler) cmdActions(ctx context.Context) string {
	items, err := b.memory.GetPendingActionItems(ctx)
	if err != nil {
		return "Failed to fetch action items."
	}
	if len(items) == 0 {
		return "No pending action items."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pending Actions (%d)\n", len(items)))
	for i, item := range items {
		if i >= 10 {
			sb.WriteString(fmt.Sprintf("... and %d more", len(items)-10))
			break
		}
		sb.WriteString(fmt.Sprintf("- %s\n", item.Description))
	}
	return sb.String()
}

func (b *BotHandler) cmdStatus(ctx context.Context) string {
	blockers, err := b.memory.GetActiveBlockers(ctx)
	if err != nil {
		b.logger.Error("get active blockers", "error", err)
	}
	risks, err := b.memory.GetOpenRisks(ctx)
	if err != nil {
		b.logger.Error("get open risks", "error", err)
	}
	actions, err := b.memory.GetPendingActionItems(ctx)
	if err != nil {
		b.logger.Error("get pending action items", "error", err)
	}
	return fmt.Sprintf("Quick Status\n- Blockers: %d\n- Risks: %d\n- Actions: %d", len(blockers), len(risks), len(actions))
}
