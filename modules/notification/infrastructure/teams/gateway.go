package teams

import (
	"context"
	"fmt"
)

// GatewayAdapter wraps the Teams client as an agent domain.Gateway.
type GatewayAdapter struct {
	client *Client
}

// NewGatewayAdapter creates a Teams gateway adapter.
func NewGatewayAdapter(client *Client) *GatewayAdapter {
	return &GatewayAdapter{client: client}
}

// Channel returns the channel name for this gateway.
func (g *GatewayAdapter) Channel() string { return "teams" }

// SendText sends a plain text message via Teams webhook.
func (g *GatewayAdapter) SendText(ctx context.Context, channelID, text string) error {
	return g.client.SendMessage(ctx, text)
}

// SendMarkdown sends a formatted Adaptive Card message via Teams webhook.
func (g *GatewayAdapter) SendMarkdown(ctx context.Context, channelID, title, content string) error {
	body := fmt.Sprintf("%s\n\n%s", title, content)
	return g.client.SendCard(ctx, title, body)
}
