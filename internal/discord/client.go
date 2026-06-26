package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/aldok10/zara-jira-mcp/config"
)

// Client wraps discordgo for messaging.
type Client struct {
	session    *discordgo.Session
	channelID  string
	webhookURL string
}

func NewClient(cfg *config.Config) *Client {
	c := &Client{
		channelID:  cfg.Discord.ChannelID,
		webhookURL: cfg.Discord.WebhookURL,
	}
	if cfg.Discord.BotToken != "" {
		s, err := discordgo.New("Bot " + cfg.Discord.BotToken)
		if err == nil {
			c.session = s
		}
	}
	return c
}

func (c *Client) Available() bool {
	return c.session != nil || c.webhookURL != ""
}

func (c *Client) SendMessage(ctx context.Context, channelID, content string) error {
	if channelID == "" {
		channelID = c.channelID
	}
	if c.session != nil {
		_, err := c.session.ChannelMessageSend(channelID, content)
		return err
	}
	return fmt.Errorf("discord bot token not configured")
}

func (c *Client) SendEmbed(ctx context.Context, channelID, title, description string, color int) error {
	if channelID == "" {
		channelID = c.channelID
	}
	if c.session == nil {
		return fmt.Errorf("discord bot token not configured")
	}
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
	_, err := c.session.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func (c *Client) ListChannels(ctx context.Context, guildID string) ([]*discordgo.Channel, error) {
	if c.session == nil {
		return nil, fmt.Errorf("discord bot token not configured")
	}
	return c.session.GuildChannels(guildID)
}
