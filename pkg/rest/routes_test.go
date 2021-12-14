package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
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

func TestMiddllewareLogging(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	sut := loggingMiddleware(observedLogger.Sugar())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest("GET", "/ping", nil)

	sut(ctx)

	logs := observedLogs.All()

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, observedLogs.Len())
	assert.NotEmpty(t, logs)
}
