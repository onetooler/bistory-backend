package main

import (
	"embed"
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
// @host localhost:8080
// @BasePath /api
func main() {
	Server{
		YamlFile:    yamlFile,
		ZapYamlFile: zapYamlFile,
		StaticFile:  staticFile,
		PropsFile:   propsFile,
	}.Run()
}
