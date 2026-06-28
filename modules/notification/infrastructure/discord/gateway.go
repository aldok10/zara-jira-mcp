package discord

import "context"

// GatewayAdapter wraps the Discord client as an agent Gateway.
type GatewayAdapter struct {
	client *Client
}

// NewGatewayAdapter creates a Discord gateway adapter.
func NewGatewayAdapter(client *Client) *GatewayAdapter {
	return &GatewayAdapter{client: client}
}

// Channel returns the channel name for this gateway.
func (g *GatewayAdapter) Channel() string { return "discord" }

// SendText sends a plain text message to a Discord channel.
func (g *GatewayAdapter) SendText(ctx context.Context, channelID, text string) error {
	return g.client.SendMessage(ctx, channelID, text)
}

// SendMarkdown sends a formatted embed message to a Discord channel.
func (g *GatewayAdapter) SendMarkdown(ctx context.Context, channelID, title, content string) error {
	return g.client.SendEmbed(ctx, channelID, title, content, 0x00AFFF)
}
