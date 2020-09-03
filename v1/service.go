package meet

import (
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/people/v1"
)

// This API requires some API scopes:
// Get user email scope
// Manage calendar event scope
//
// Note that this scope is slice of string, not string
var MeetScopes = []string{
	people.UserinfoEmailScope,
	calendar.CalendarEventsScope,
}

// Main API service struct
type Service struct {
	Conference *ConferenceService
}

// Create server pointer
func New(c *http.Client) (*Service, error) {
	cs, err := NewConferenceService(c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize ConferenceCalService")
	}
	return &Service{
		Conference: cs,
	}, nil
}
