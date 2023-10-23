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
	"github.com/onetooler/bistory-backend/router"
	"github.com/onetooler/bistory-backend/session"
)

//go:embed resources/config/application.*.yml
var yamlFile embed.FS

//go:embed resources/config/zaplogger.*.yml
var zapYamlFile embed.FS

//go:embed resources/public/*
var staticFile embed.FS

//go:embed resources/config/messages.properties
var propsFile embed.FS

// @title bistory-backend API
// @version 0.0.1
// @description This is API specification for bistory-backend project.

// @license.name MIT
// @license.url https://opensource.org/licenses/mit-license.php

// @host localhost:8080
// @BasePath /api
func main() {
	e := echo.New()

	conf, env := config.LoadAppConfig(yamlFile)
	logger := logger.InitLogger(env, zapYamlFile)
	logger.GetZapLogger().Infof("Loaded this configuration : application." + env + ".yml")

	messages := config.LoadMessagesConfig(propsFile)
	logger.GetZapLogger().Infof("Loaded messages.properties")

	rep := repository.NewBookRepository(logger, conf)
	sess := session.NewSession()
	container := container.NewContainer(rep, sess, conf, messages, logger, env)

	migration.CreateDatabase(container)
	migration.InitMasterData(container)

	router.Init(e, container)
	middleware.InitLoggerMiddleware(e, container)
	middleware.InitSessionMiddleware(e, container)
	middleware.StaticContentsMiddleware(e, container, staticFile)

	if err := e.Start(":8080"); err != nil {
		logger.GetZapLogger().Errorf(err.Error())
	}

	defer rep.Close()
}
