package rest

import (
	"bytes"
	"io"
	"net/http"

	"example.com/playground/pkg/api"
	"example.com/playground/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

func SetupRouteHandlers(r *RouteHandlers, lg *zap.SugaredLogger) *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(loggingMiddleware(lg))
	e.SetTrustedProxies(nil)
	return r.routes(e)
}

func loggingMiddleware(lg *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b []byte
		if c.Request.Body != nil {
			buf := new(bytes.Buffer)

			if _, err := buf.ReadFrom(c.Request.Body); err != nil {
				lg.Errorf("serialization error", "error", err)

				return
			}

			b = buf.Bytes()

			c.Request.Body = io.NopCloser(buf)
		}

		c.Next()

		lg.Infow("request",
			"requestURI", c.Request.RequestURI,
			"header", c.Request.Header,
			"host", c.Request.Host,
			"method", c.Request.Method,
			"body", string(b))
	}
}

type RouteHandlers struct {
	CreateOffer storage.CreateOffer
	UpdateOffer storage.UpdateOffer
	GetOffer    storage.GetOffer
	DeleteOffer storage.DeleteOffer
}

func (r *RouteHandlers) routes(e *gin.Engine) *gin.Engine {

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	e.POST("/offers", create(r.CreateOffer))
	e.PUT("/offers/:offerID", update(r.UpdateOffer))
	e.GET("/offers/:offerID", getByID(r.GetOffer))
	e.DELETE("/offers/:offerID", delete(r.DeleteOffer))

	return e
}

func delete(d storage.DeleteOffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Header("Content-Type", "application/json")

		offerID := c.Param("offerID")

		if err := validation.Validate(offerID,
			validation.Required,
			is.UUID,
		); err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)

			return
		}

		err := d.DeleteByID(c.Request.Context(), offerID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.Status(http.StatusOK)

	}
}

func create(cr storage.CreateOffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Header("Content-Type", "application/json")

		var request api.JobOfferRequest
		if err := c.BindJSON(&request); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		resp, err := cr.Create(c.Request.Context(), &request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func update(u storage.UpdateOffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Header("Content-Type", "application/json")

		var request api.UpdateJobOfferRequest
		if err := c.BindJSON(&request); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		offerID := c.Param("offerID")

		if err := validation.Validate(offerID,
			validation.Required,
			is.UUID,
		); err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)

			return
		}

		resp, err := u.Update(c.Request.Context(), offerID, &request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func getByID(g storage.GetOffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Header("Content-Type", "application/json")

		offerID := c.Param("offerID")

		if err := validation.Validate(offerID,
			validation.Required,
			is.UUID,
		); err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)

			return
		}

		resp, err := g.Get(c.Request.Context(), offerID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
