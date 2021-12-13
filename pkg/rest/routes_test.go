package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestRoutes(t *testing.T) {
	sut := SetupRouteHandlers(&RouteHandlers{}, zaptest.NewLogger(t).Sugar())

	routes := make(map[string][]string)
	for _, r := range sut.Routes() {
		if routes[r.Path] == nil {
			routes[r.Path] = []string{}
		}

		routes[r.Path] = append(routes[r.Path], r.Method)
	}

	assert.Equal(t, map[string][]string{
		"/ping":            {"GET"},
		"/offers":          {"POST"},
		"/offers/:offerID": {"GET", "PUT", "DELETE"},
	}, routes)

}
