package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/test"
	"github.com/onetooler/bistory-backend/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogging(t *testing.T) {
	router, container, logs := test.PrepareForLoggerTest()

	book := NewBookController(container)
	router.GET(config.APIBooksID, func(c echo.Context) error { return book.GetBook(c) })

	setUpTestData(container)

	uri := util.NewRequestBuilder().URL(config.APIBooks).PathParams("1").Build().GetRequestURL()
	req := httptest.NewRequest("GET", uri, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	allLogs := logs.All()
	assert.True(t, assertLogger("/api/books/:id Action Start", allLogs))
	assert.True(t, assertLogger("/api/books/:id Action End", allLogs))
	assert.True(t, assertLogger("/api/books/1 GET 200", allLogs))
	assert.True(t, assertLogger("[gorm] ", allLogs))
}

func assertLogger(message string, logs []observer.LoggedEntry) bool {
	for _, l := range logs {
		if strings.Contains(l.Message, message) {
			return true
		}
	}
	return false
}
