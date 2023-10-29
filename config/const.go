package config

// ErrExitStatus represents the error status in this application.
const ErrExitStatus int = 2

const (
	// AppConfigPath is the path of application.yml.
	AppConfigPath = "resources/config/application.%s.yml"
	// MessagesConfigPath is the path of messages.properties.
	MessagesConfigPath = "resources/config/messages.properties"
	// LoggerConfigPath is the path of zaplogger.yml.
	LoggerConfigPath = "resources/config/zaplogger.%s.yml"
)

// PasswordHashCost is hash cost for a password.
const PasswordHashCost int = 10

const (
	// API represents the group of API.
	API = "/api"
)

const (
	// APIAuth represents the group of auth management API.
	// APIAuth is a part of account but separate API for concentrate scope
	APIAuth = API + "/auth"
	APIAuthLoginStatus = APIAuth + "/loginStatus"
	APIAuthLoginAccount = APIAuth + "/loginAccount"
	APIAuthLogin = APIAuth + "/login"
	APIAuthLogout = APIAuth + "/logout"
)

const (
	// APIAccount represents the group of account management API.
	APIAccount = API + "/account"
	APIAccountIdParam = "id"
	APIAccountIdPath = APIAccount + "/:" + APIAccountIdParam
)

const (
	// APIHealth represents the API to get the status of this application.
	APIHealth = API + "/health"
)
