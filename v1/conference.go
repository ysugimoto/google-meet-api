package meet

import (
	"context"
	"strings"
	"time"

	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/api/calendar/v3"
	user "google.golang.org/api/oauth2/v2"
)

const (
	primaryCalendarId = "primary"
	googleMeetDomain  = "https://meet.google.com"
)

// Conference is struct for api call response data
type Conference struct {
	URL string `json:"url,omitempty"`
}

// ConferenceService wraps calendar's conferenceData access
type ConferenceService struct {
	user     *user.Service
	calendar *calendar.Service
}

// Create new ConferenceService
func NewConferenceService(c *http.Client) (*ConferenceService, error) {
	u, err := user.New(c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize user API client")
	}
	cal, err := calendar.New(c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize calendar API client")
	}
	return &ConferenceService{
		user:     u,
		calendar: cal,
	}, nil
}

// Create() creates Calendar creation caller pointer
func (c *ConferenceService) Create() *ConferenceCreateCall {
	return &ConferenceCreateCall{
		user:     c.user,
		calendar: c.calendar,
	}
}

// ConferenceCreateCall manages creation of ConferenceData
type ConferenceCreateCall struct {
	user     *user.Service
	calendar *calendar.Service
	ctx      context.Context
}

// Set request context
func (c *ConferenceCreateCall) Context(ctx context.Context) *ConferenceCreateCall {
	c.ctx = ctx
	return c
}

// Get request context
func (c *ConferenceCreateCall) context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

// Execute API
func (c *ConferenceCreateCall) Do() (*Conference, error) {
	ctx := c.context()

	// Get user email to use default conference attendee
	u, err := c.user.Userinfo.Get().Context(ctx).Fields("email").Do()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authed user info")
	}

	event := &calendar.Event{
		Summary: "Temporal Event (will be deleted immediately)",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
		},
		Attendees: []*calendar.EventAttendee{
			{
				Email: u.Email,
			},
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: uuid.New().String(),
			},
		},
	}

	ret, err := c.calendar.Events.Insert(primaryCalendarId, event).
		ConferenceDataVersion(1).
		Context(ctx).
		Do()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to insert temporary event")
	}
	// make sure temporal event should be deleted
	defer c.calendar.Events.Delete(primaryCalendarId, ret.Id).Do()

	var meetURL string
	if ret.ConferenceData != nil {
		for _, entry := range ret.ConferenceData.EntryPoints {
			if entry.EntryPointType != "video" {
				continue
			}
			if strings.HasPrefix(entry.Uri, googleMeetDomain) {
				meetURL = entry.Uri
			}
		}
	}
	if meetURL == "" {
		return nil, errors.New("Failed to retrieve google meet url: not found in response")
	}

	return &Conference{
		URL: meetURL,
	}, nil
}
