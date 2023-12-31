package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthCheck(t *testing.T) {
	router, container := testutil.PrepareForControllerTest(false)

	health := NewHealthController(container)
	router.GET(config.APIHealth, func(c echo.Context) error { return health.GetHealthCheck(c) })

	req := httptest.NewRequest("GET", config.APIHealth, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `"healthy"`, rec.Body.String())
}
