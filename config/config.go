package config

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"

	"github.com/onetooler/bistory-backend/util"
	"gopkg.in/yaml.v3"
)

// Config represents the composition of yml settings.
type Config struct {
	Database struct {
		Dialect   string `default:"sqlite3"`
		Host      string `default:"develop.db"`
		Port      string
		Dbname    string
		LoginId   string
		Password  string
		Migration bool `default:"false"`
	}
	Redis struct {
		Enabled            bool `default:"false"`
		ConnectionPoolSize int  `yaml:"connection_pool_size" default:"10"`
		Host               string
		Port               string
	}
	Email struct {
		Account  string
		Host     string
		Port     int
		Username string
		Password string
	}
	Extension struct {
		MasterGenerator bool `yaml:"master_generator" default:"false"`
		CorsEnabled     bool `yaml:"cors_enabled" default:"false"`
		SecurityEnabled bool `yaml:"security_enabled" default:"false"`
	}
	Log struct {
		RequestLogFormat string `yaml:"request_log_format" default:"${remote_ip} ${account_loginid} ${uri} ${method} ${status}"`
	}
	StaticContents struct {
		Enabled bool `default:"false"`
	}
	Swagger struct {
		Enabled bool `default:"false"`
		Path    string
	}
	Security struct {
		AuthPath    []string `yaml:"auth_path"`
		ExcludePath []string `yaml:"exclude_path"`
		UserPath    []string `yaml:"user_path"`
		AdminPath   []string `yaml:"admin_path"`
	}
}

const (
	// DEV represents development environment
	DEV = "develop"
	// PRD represents production environment
	PRD = "production"
	// DOC represents docker container
	DOC = "docker"
)

// LoadAppConfig reads the settings written to the yml file
func LoadAppConfig(yamlFile embed.FS) (*Config, string) {
	var env *string
	if value := os.Getenv("WEB_APP_ENV"); value != "" {
		env = &value
	} else {
		env = flag.String("env", "develop", "To switch configurations.")
		flag.Parse()
	}

	file, err := yamlFile.ReadFile(fmt.Sprintf(AppConfigPath, *env))
	if err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	return config, *env
}

// LoadMessagesConfig loads the messages.properties.
func LoadMessagesConfig(propsFile embed.FS) map[string]string {
	messages := util.ReadPropertiesFile(propsFile, MessagesConfigPath)
	if messages == nil {
		fmt.Printf("Failed to load the messages.properties.")
		os.Exit(ErrExitStatus)
	}
	return messages
}

func LoadEmailTemplates(emailFile embed.FS) map[string]*template.Template {
	subfs, err := fs.Sub(emailFile, EmailTemplatesPath)
	if err != nil {
		fmt.Printf("Failed to load the email templates.")
		os.Exit(ErrExitStatus)
	}

	files, err := util.GetAllFiles(subfs)
	if err != nil {
		fmt.Printf("Failed to load the email templates.")
		os.Exit(ErrExitStatus)
	}

	templates := make(map[string]*template.Template)
	for fileName, fileBody := range files {
		t, err := template.New(fileName).Parse(fileBody)
		if err != nil {
			fmt.Printf("Failed to load the email templates.")
			os.Exit(ErrExitStatus)
		}
		templates[fileName] = t
	}

	return templates
}
