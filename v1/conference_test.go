package meet_test

import (
	"strings"
	"testing"

	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/ysugimoto/google-meet-api/v1"
)

const eventFixture = `
{
	"id": "test",
	"conferenceData": {
		"entryPoints": [
			{
				"entryPointType": "video",
				"uri": "https://meet.google.com/test-meet-uri"
			}
		]
	}
}`

const userFixture = `
{
	"email": "test@example.com"
}`

type RoungTripper struct{}

func (r *RoungTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	switch {
	case strings.HasPrefix(req.URL.Path, "/calendar/v3"):
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(eventFixture)),
		}, nil
	case strings.HasPrefix(req.URL.Path, "/oauth2/v2/userinfo"):
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(userFixture)),
		}, nil
	default:
		return &http.Response{
			StatusCode: http.StatusNotFound,
		}, nil
	}
}

var mockClient = &http.Client{
	Transport: &RoungTripper{},
}

func TestConferenceService(t *testing.T) {
	srv, err := meet.NewConferenceService(mockClient)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	ret, err := srv.Create().Do()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "https://meet.google.com/test-meet-uri", ret.URL)
}
