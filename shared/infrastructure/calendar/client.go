package calendar

import (
	"context"
	"fmt"
	"strconv"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcal "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// Client provides calendar operations via Lark and/or Google Calendar.
type Client struct {
	lark       *lark.Client
	calendarID string // Lark primary calendar ID (auto-fetched if empty)
}

func NewClient(cfg *config.Config) *Client {
	c := &Client{}
	if cfg.Lark.AppID != "" && cfg.Lark.AppSecret != "" {
		c.lark = lark.NewClient(cfg.Lark.AppID, cfg.Lark.AppSecret)
	}
	return c
}

func (c *Client) Available() bool {
	return c.lark != nil
}

// Event represents a calendar event.
type Event struct {
	ID          string
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	Attendees   []string
	MeetingURL  string
}

// CreateEvent creates a calendar event in Lark.
func (c *Client) CreateEvent(ctx context.Context, ev *Event) (*Event, error) {
	if c.lark == nil {
		return nil, fmt.Errorf("calendar not configured: set LARK_APP_ID + LARK_APP_SECRET")
	}

	calID, err := c.getPrimaryCalendarID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get calendar: %w", err)
	}

	startTS := strconv.FormatInt(ev.StartTime.Unix(), 10)
	endTS := strconv.FormatInt(ev.EndTime.Unix(), 10)
	tz := "Asia/Jakarta"

	req := larkcal.NewCreateCalendarEventReqBuilder().
		CalendarId(calID).
		CalendarEvent(larkcal.NewCalendarEventBuilder().
			Summary(ev.Summary).
			Description(ev.Description).
			StartTime(larkcal.NewTimeInfoBuilder().Timestamp(startTS).Timezone(tz).Build()).
			EndTime(larkcal.NewTimeInfoBuilder().Timestamp(endTS).Timezone(tz).Build()).
			Build()).
		Build()

	resp, err := c.lark.Calendar.CalendarEvent.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, fmt.Errorf("lark calendar error %d: %s", resp.Code, resp.Msg)
	}

	created := &Event{
		ID:      *resp.Data.Event.EventId,
		Summary: ev.Summary,
	}
	return created, nil
}

// ListEvents lists upcoming events.
func (c *Client) ListEvents(ctx context.Context, from, to time.Time) ([]Event, error) {
	if c.lark == nil {
		return nil, fmt.Errorf("calendar not configured")
	}

	calID, err := c.getPrimaryCalendarID(ctx)
	if err != nil {
		return nil, err
	}

	startTS := strconv.FormatInt(from.Unix(), 10)
	endTS := strconv.FormatInt(to.Unix(), 10)

	req := larkcal.NewListCalendarEventReqBuilder().
		CalendarId(calID).
		StartTime(startTS).
		EndTime(endTS).
		Build()

	resp, err := c.lark.Calendar.CalendarEvent.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, fmt.Errorf("lark calendar error %d: %s", resp.Code, resp.Msg)
	}

	events := make([]Event, 0, len(resp.Data.Items))
	for _, item := range resp.Data.Items {
		ev := Event{
			ID: safeStr(item.EventId),
		}
		if item.Summary != nil {
			ev.Summary = *item.Summary
		}
		if item.Description != nil {
			ev.Description = *item.Description
		}
		if item.StartTime != nil && item.StartTime.Timestamp != nil {
			ts, _ := strconv.ParseInt(*item.StartTime.Timestamp, 10, 64)
			ev.StartTime = time.Unix(ts, 0)
		}
		if item.EndTime != nil && item.EndTime.Timestamp != nil {
			ts, _ := strconv.ParseInt(*item.EndTime.Timestamp, 10, 64)
			ev.EndTime = time.Unix(ts, 0)
		}
		if item.Location != nil && item.Location.Name != nil {
			ev.Location = *item.Location.Name
		}
		events = append(events, ev)
	}
	return events, nil
}

// ScheduleMeeting creates a meeting event with video conference link.
func (c *Client) ScheduleMeeting(ctx context.Context, summary, description string, start, end time.Time, attendeeEmails []string) (*Event, error) {
	if c.lark == nil {
		return nil, fmt.Errorf("calendar not configured")
	}

	calID, err := c.getPrimaryCalendarID(ctx)
	if err != nil {
		return nil, err
	}

	startTS := strconv.FormatInt(start.Unix(), 10)
	endTS := strconv.FormatInt(end.Unix(), 10)
	tz := "Asia/Jakarta"
	vcType := "vc" // Lark video conference

	req := larkcal.NewCreateCalendarEventReqBuilder().
		CalendarId(calID).
		CalendarEvent(larkcal.NewCalendarEventBuilder().
			Summary(summary).
			Description(description).
			StartTime(larkcal.NewTimeInfoBuilder().Timestamp(startTS).Timezone(tz).Build()).
			EndTime(larkcal.NewTimeInfoBuilder().Timestamp(endTS).Timezone(tz).Build()).
			Vchat(larkcal.NewVchatBuilder().VcType(vcType).Build()).
			Build()).
		Build()

	resp, err := c.lark.Calendar.CalendarEvent.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, fmt.Errorf("lark calendar error %d: %s", resp.Code, resp.Msg)
	}

	created := &Event{
		ID:      *resp.Data.Event.EventId,
		Summary: summary,
	}

	// Add attendees if provided
	if len(attendeeEmails) > 0 {
		_ = c.addAttendees(ctx, calID, *resp.Data.Event.EventId, attendeeEmails)
	}

	return created, nil
}

func (c *Client) addAttendees(ctx context.Context, calendarID, eventID string, emails []string) error {
	attendees := make([]*larkcal.CalendarEventAttendee, 0, len(emails))
	for _, email := range emails {
		attendees = append(attendees, larkcal.NewCalendarEventAttendeeBuilder().
			Type("third_party").
			ThirdPartyEmail(email).
			Build())
	}

	req := larkcal.NewCreateCalendarEventAttendeeReqBuilder().
		CalendarId(calendarID).
		EventId(eventID).
		Body(larkcal.NewCreateCalendarEventAttendeeReqBodyBuilder().
			Attendees(attendees).
			Build()).
		Build()

	resp, err := c.lark.Calendar.CalendarEventAttendee.Create(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("add attendees error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}

func (c *Client) getPrimaryCalendarID(ctx context.Context) (string, error) {
	if c.calendarID != "" {
		return c.calendarID, nil
	}

	req := larkcal.NewPrimaryCalendarReqBuilder().Build()
	resp, err := c.lark.Calendar.Calendar.Primary(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("get primary calendar error %d: %s", resp.Code, resp.Msg)
	}

	c.calendarID = *resp.Data.Calendars[0].Calendar.CalendarId
	return c.calendarID, nil
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
