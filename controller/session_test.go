package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	_ "github.com/onetooler/bistory-backend/docs" // for using echo-swagger
	"github.com/onetooler/bistory-backend/testutil"
	"github.com/stretchr/testify/assert"
)

type sessionController struct {
	container container.Container
}

func TestSessionRace_Success(t *testing.T) {
	sessionKey := "Key"
	router, container := testutil.PrepareForControllerTest(true, false)
	session := sessionController{container: container}
	router.GET(config.API+"1", func(c echo.Context) error {
		session.container.GetSession().SetValue(sessionKey, 1)
		session.container.GetSession().Save()
		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, session.container.GetSession().GetValue(sessionKey))
	})
	router.GET(config.API+"2", func(c echo.Context) error {
		session.container.GetSession().SetValue(sessionKey, 2)
		session.container.GetSession().Save()
		return c.String(http.StatusOK, session.container.GetSession().GetValue(sessionKey))
	})

	req1 := httptest.NewRequest("GET", config.API+"1", nil)
	req2 := httptest.NewRequest("GET", config.API+"2", nil)
	rec1 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()

	go func() {
		router.ServeHTTP(rec1, req1)
	}()

	go func() {
		time.Sleep(1 * time.Second)
		router.ServeHTTP(rec2, req2)
	}()

	time.Sleep(5 * time.Second)

	assert.Equal(t, "1", rec1.Body.String())
	assert.Equal(t, "2", rec2.Body.String())
}
