package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/playground/pkg/api"

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
		"/offers":          {"GET", "POST"},
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

type testDeleteOffer struct {
	deleteOfferCalled int
	deleteOfferErr    error
}

func (d *testDeleteOffer) DeleteByID(ctx context.Context, offerID string) error {
	d.deleteOfferCalled++
	return d.deleteOfferErr
}

func TestDeleteOffer(t *testing.T) {
	tests := []struct {
		offerID        string
		expectedCalls  int
		expectedStatus int
		expectedErr    error
	}{
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			1,
			http.StatusOK,
			nil,
		},
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			1,
			http.StatusInternalServerError,
			errors.New("oops..something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			w := httptest.NewRecorder()

			do := testDeleteOffer{
				deleteOfferErr: test.expectedErr,
			}

			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = []gin.Param{
				{
					Key:   "offerID",
					Value: test.offerID,
				},
			}
			ctx.Request = httptest.NewRequest("DELETE", fmt.Sprintf("/offers/%s", test.offerID), nil)

			delete(&do)(ctx)

			assert.Equal(t, test.expectedStatus, ctx.Writer.Status())
			assert.Equal(t, test.expectedCalls, do.deleteOfferCalled)
		})
	}
}

type testCreateOffer struct {
	createOfferCalled int
	createOfferErr    error
}

func (t *testCreateOffer) Create(context.Context, *api.JobOfferRequest) (*api.JobOfferResponse, error) {
	t.createOfferCalled++
	return &api.JobOfferResponse{}, t.createOfferErr
}

func TestCreateOffer(t *testing.T) {
	tests := []struct {
		request        *api.JobOfferRequest
		expectedCalls  int
		expectedStatus int
		expectedErr    error
	}{
		{
			&api.JobOfferRequest{
				Company:        "TEST",
				Email:          "test@hr-test.com",
				ExpirationDate: "2022-03-01 14:30:00.00000",
				LinkToOffer:    "http://test.com/carriers",
				Details:        "We are looking for a Ninja Golang developer to work on our system serving...",
				Salary:         18000,
				ContactPhone:   "+38978653534",
			},
			1,
			http.StatusCreated,
			nil,
		},
		{
			&api.JobOfferRequest{},
			1,
			http.StatusInternalServerError,
			errors.New("oops..something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			w := httptest.NewRecorder()

			co := testCreateOffer{
				createOfferErr: test.expectedErr,
			}

			ctx, _ := gin.CreateTestContext(w)

			b, err := json.Marshal(test.request)

			assert.Nil(t, err)

			ctx.Request = httptest.NewRequest("POST", "/offers", bytes.NewReader(b))

			create(&co)(ctx)

			assert.Equal(t, test.expectedStatus, ctx.Writer.Status())
			assert.Equal(t, test.expectedCalls, co.createOfferCalled)
		})
	}
}

type testUpdateOffer struct {
	updateOfferCalled int
	updateOfferErr    error
}

func (t *testUpdateOffer) Update(context.Context, string, *api.UpdateJobOfferRequest) (*api.JobOfferResponse, error) {
	t.updateOfferCalled++
	return &api.JobOfferResponse{}, t.updateOfferErr
}

func TestUpdateOffer(t *testing.T) {
	tests := []struct {
		offerID        string
		request        *api.UpdateJobOfferRequest
		expectedCalls  int
		expectedStatus int
		expectedErr    error
	}{
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			&api.UpdateJobOfferRequest{
				Salary:       8500,
				Email:        "hey@outlook.com",
				ContactPhone: "+38978360298",
				LinkToOffer:  "http://test.com/carrers",
			},
			1,
			http.StatusOK,
			nil,
		},
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			&api.UpdateJobOfferRequest{},
			1,
			http.StatusInternalServerError,
			errors.New("oops..something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			w := httptest.NewRecorder()

			u := testUpdateOffer{
				updateOfferErr: test.expectedErr,
			}

			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = []gin.Param{
				{
					Key:   "offerID",
					Value: test.offerID,
				},
			}
			b, err := json.Marshal(test.request)

			assert.Nil(t, err)

			ctx.Request = httptest.NewRequest("PUT", fmt.Sprintf("/offers/%s", test.offerID), bytes.NewReader(b))

			update(&u)(ctx)

			assert.Equal(t, test.expectedStatus, ctx.Writer.Status())
			assert.Equal(t, test.expectedCalls, u.updateOfferCalled)
		})
	}
}

type testGetOffer struct {
	getOfferCalled int
	getOfferErr    error
}

func (d *testGetOffer) Get(context.Context, string) (*api.JobOfferResponse, error) {
	d.getOfferCalled++
	return &api.JobOfferResponse{}, d.getOfferErr
}

func TestGetOfferByID(t *testing.T) {
	tests := []struct {
		offerID        string
		expectedCalls  int
		expectedStatus int
		expectedErr    error
	}{
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			1,
			http.StatusOK,
			nil,
		},
		{
			"eca51142-3bf0-4766-baf7-2a168c964024",
			1,
			http.StatusInternalServerError,
			errors.New("oops..something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			w := httptest.NewRecorder()

			g := testGetOffer{
				getOfferErr: test.expectedErr,
			}

			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = []gin.Param{
				{
					Key:   "offerID",
					Value: test.offerID,
				},
			}
			ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/offers/%s", test.offerID), nil)

			getByID(&g)(ctx)

			assert.Equal(t, test.expectedStatus, ctx.Writer.Status())
			assert.Equal(t, test.expectedCalls, g.getOfferCalled)
		})
	}
}
