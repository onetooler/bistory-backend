package main

import (
	"embed"

	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/container"
	"github.com/onetooler/bistory-backend/logger"
	"github.com/onetooler/bistory-backend/middleware"
	"github.com/onetooler/bistory-backend/migration"
	"github.com/onetooler/bistory-backend/repository"
	"github.com/onetooler/bistory-backend/routes"
	"github.com/onetooler/bistory-backend/session"
	"github.com/onetooler/bistory-backend/util"
)

type Server struct {
	YamlFile    embed.FS
	ZapYamlFile embed.FS
	StaticFile  embed.FS
	PropsFile   embed.FS
}

func (s Server) Run() {
	e := echo.New()

	conf, env := config.LoadAppConfig(s.YamlFile)
	logger := logger.InitLogger(env, s.ZapYamlFile)
	logger.GetZapLogger().Infof("Loaded this configuration : application." + env + ".yml")

	messages := config.LoadMessagesConfig(s.PropsFile)
	logger.GetZapLogger().Infof("Loaded messages.properties")

	rep := repository.NewRepository(logger, conf)
	defer util.Check(rep.Close)

	sess := session.NewSession()
	container := container.NewContainer(rep, sess, conf, messages, logger, env)

	migration.Init(container)
	routes.Init(e, container)
	middleware.Init(e, container, s.StaticFile)

	if err := e.Start(":8080"); err != nil {
		logger.GetZapLogger().Errorf(err.Error())
	}
}
