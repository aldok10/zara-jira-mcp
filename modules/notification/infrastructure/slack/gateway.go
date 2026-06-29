package slack

import "context"

// GatewayAdapter wraps the Slack client as an agent Gateway.
type GatewayAdapter struct {
	client *Client
}

// NewGatewayAdapter creates a Slack gateway adapter.
func NewGatewayAdapter(client *Client) *GatewayAdapter {
	return &GatewayAdapter{client: client}
}

// Channel returns the channel name for this gateway.
func (g *GatewayAdapter) Channel() string { return "slack" }

// SendText sends a plain text message to a Slack channel.
func (g *GatewayAdapter) SendText(ctx context.Context, channelID, text string) error {
	return g.client.SendMessage(ctx, channelID, text)
}

// SendMarkdown sends a formatted message with title to a Slack channel.
func (g *GatewayAdapter) SendMarkdown(ctx context.Context, channelID, title, content string) error {
	return g.client.SendRichMessage(ctx, channelID, title, content)
}
