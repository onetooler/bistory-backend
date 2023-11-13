package middleware

import (
	"embed"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomd "github.com/labstack/echo/v4/middleware"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/model"
	"github.com/valyala/fasttemplate"
)

var authorizationPathRegexps map[string]*regexp.Regexp

func Init(e *echo.Echo, container container.Container, staticFile embed.FS) {
	InitCORSMiddleware(e, container)
	InitLoggerMiddleware(e, container)
	InitSessionMiddleware(e, container)
	StaticContentsMiddleware(e, container, staticFile)
}

func InitCORSMiddleware(e *echo.Echo, container container.Container) {
	if container.GetConfig().Extension.CorsEnabled {
		e.Use(echomd.CORSWithConfig(echomd.CORSConfig{
			AllowCredentials:                         true,
			UnsafeWildcardOriginWithAllowCredentials: true,
			AllowOrigins:                             []string{"*"},
			AllowHeaders: []string{
				echo.HeaderAccessControlAllowHeaders,
				echo.HeaderContentType,
				echo.HeaderContentLength,
				echo.HeaderAcceptEncoding,
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			MaxAge: 86400,
		}))
	}
}

// InitLoggerMiddleware initialize a middleware for logger.
func InitLoggerMiddleware(e *echo.Echo, container container.Container) {
	e.Use(RequestLoggerMiddleware(container))
	e.Use(ActionLoggerMiddleware(container))
	e.Use(BodyLoggerMiddleware(container))
}

// InitSessionMiddleware initialize a middleware for session management.
func InitSessionMiddleware(e *echo.Echo, container container.Container) {
	e.Use(session.Middleware(container.GetSession().GetStore()))
	e.Use(AuthenticationMiddleware(container))
}

// RequestLoggerMiddleware is middleware for logging the contents of requests.
func RequestLoggerMiddleware(container container.Container) echo.MiddlewareFunc {
	template := fasttemplate.New(container.GetConfig().Log.RequestLogFormat, "${", "}")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			if err := next(c); err != nil {
				c.Error(err)
			}

			logstr := template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "remote_ip":
					return w.Write([]byte(c.RealIP()))
				case "account_loginid":
					if account := container.GetSession().GetAccount(c); account != nil {
						return w.Write([]byte(account.LoginId))
					}
					return w.Write([]byte("None"))
				case "uri":
					return w.Write([]byte(req.RequestURI))
				case "method":
					return w.Write([]byte(req.Method))
				case "status":
					return w.Write([]byte(strconv.Itoa(res.Status)))
				default:
					return w.Write([]byte(""))
				}
			})
			container.GetLogger().GetZapLogger().Infof(logstr)
			return nil
		}
	}
}

// ActionLoggerMiddleware is middleware for logging the start and end of controller processes.
// ref: https://echo.labstack.com/cookbook/middleware
func ActionLoggerMiddleware(container container.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := container.GetLogger()
			logger.GetZapLogger().Debugf("%s Action Start", c.Path())
			if err := next(c); err != nil {
				c.Error(err)
			}
			logger.GetZapLogger().Debugf("%s Action End", c.Path())
			return nil
		}
	}
}

func BodyLoggerMiddleware(container container.Container) echo.MiddlewareFunc {
	return echomd.BodyDumpWithConfig(
		echomd.BodyDumpConfig{
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "swagger")
			},
			Handler: func(c echo.Context, reqBody []byte, resBody []byte) {
				logger := container.GetLogger()
				logger.GetZapLogger().Debugf("%s request body: %s", c.Path(), reqBody)
				logger.GetZapLogger().Debugf("%s response body: %s", c.Path(), resBody)
			},
		},
	)
}

// StaticContentsMiddleware is the middleware for loading the static files.
func StaticContentsMiddleware(e *echo.Echo, container container.Container, staticFile embed.FS) {
	conf := container.GetConfig()
	if conf.StaticContents.Enabled {
		staticConfig := echomd.StaticConfig{
			Root:       "resources/public",
			Index:      "index.html",
			Browse:     false,
			HTML5:      true,
			Filesystem: http.FS(staticFile),
		}
		if conf.Swagger.Enabled {
			staticConfig.Skipper = func(c echo.Context) bool {
				return equalPath(c.Path(), []string{conf.Swagger.Path})
			}
		}
		e.Use(echomd.StaticWithConfig(staticConfig))
		container.GetLogger().GetZapLogger().Infof("Served the static contents.")
	}
}

// AuthenticationMiddleware is the middleware of session authentication for echo.
func AuthenticationMiddleware(container container.Container) echo.MiddlewareFunc {
	// pre-compile all paths
	authorizationPathRegexps = make(map[string]*regexp.Regexp)
	for _, path := range container.GetConfig().Security.AuthPath {
		authorizationPathRegexps[path] = regexp.MustCompile(path)
	}
	for _, path := range container.GetConfig().Security.ExcludePath {
		authorizationPathRegexps[path] = regexp.MustCompile(path)
	}
	for _, path := range container.GetConfig().Security.AdminPath {
		authorizationPathRegexps[path] = regexp.MustCompile(path)
	}
	for _, path := range container.GetConfig().Security.UserPath {
		authorizationPathRegexps[path] = regexp.MustCompile(path)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !hasAuthorization(c, container) {
				return c.JSON(http.StatusForbidden, false)
			}
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// hasAuthorization judges whether the user has the right to access the path.
func hasAuthorization(c echo.Context, container container.Container) bool {
	currentPath := c.Path()
	if !equalPath(currentPath, container.GetConfig().Security.AuthPath) {
		return true
	}
	if equalPath(currentPath, container.GetConfig().Security.ExcludePath) {
		return true
	}

	account := container.GetSession().GetAccount(c)
	if account == nil {
		return false
	}
	if account.Authority == uint(model.AuthorityAdmin) && equalPath(currentPath, container.GetConfig().Security.AdminPath) {
		_ = container.GetSession().Save(c)
		return true
	}
	if account.Authority <= uint(model.AuthorityUser) && equalPath(currentPath, container.GetConfig().Security.UserPath) {
		_ = container.GetSession().Save(c)
		return true
	}

	return false
}

// equalPath judges whether a given path contains in the path list.
func equalPath(cpath string, paths []string) bool {
	for _, path := range paths {
		if authorizationPathRegexps[path].Match([]byte(cpath)) {
			return true
		}
	}
	return false
}
