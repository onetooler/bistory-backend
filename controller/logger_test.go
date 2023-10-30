package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogging(t *testing.T) {
	router, container, logs := test.PrepareForLoggerTest()

	health := NewHealthController(container)
	router.GET(config.APIHealth, func(c echo.Context) error { return health.GetHealthCheck(c) })

	req := httptest.NewRequest("GET", config.APIHealth, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	allLogs := logs.All()
	assert.True(t, assertLogger("/api/health Action Start", allLogs))
	assert.True(t, assertLogger("/api/health Action End", allLogs))
	assert.True(t, assertLogger("None /api/health GET 200", allLogs))
}

func assertLogger(message string, logs []observer.LoggedEntry) bool {
	for _, l := range logs {
		if strings.Contains(l.Message, message) {
			return true
		}
	}
	return false
}
