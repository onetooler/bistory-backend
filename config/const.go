package config

// ErrExitStatus represents the error status in this application.
const ErrExitStatus int = 2

const (
	EmailTemplatesPath        = "resources/email"
	FindLoginIdTemplate       = "find-login-id.html"
	EmailVerificationTemplate = "email-verification.html"

	AppConfigPath      = "resources/config/application.%s.yml"
	MessagesConfigPath = "resources/config/messages.properties"
	LoggerConfigPath   = "resources/config/zaplogger.%s.yml"
)

// Constant about account&auth domain
const (
	PasswordHashCost             int = 10
	MaxLoginAttempts             int = 5
	EmailVerificationTokenLength int = 6
)

const (
	// API represents the group of API.
	API = "/api"
)

const (
	// APIAuth represents the group of auth management API.
	// APIAuth is a part of account but separate API for concentrate scope
	APIAuth             = API + "/auth"
	APIAuthLoginStatus  = APIAuth + "/loginStatus"
	APIAuthLoginAccount = APIAuth + "/loginAccount"
	APIAuthLogin        = APIAuth + "/login"
	APIAuthLogout       = APIAuth + "/logout"

	APIAuthEmailVerificationTokenSend = APIAuth + "/email-verification"
)

const (
	// APIAccount represents the group of account management API.
	APIAccount               = API + "/account"
	APIAccountFindLoginId    = APIAccount + "/find-login-id"
	APIAccountIdParam        = "id"
	APIAccountLoginIdParam   = "loginid"
	APIAccountIdPath         = APIAccount + "/:" + APIAccountIdParam
	APIAccountChangePassword = APIAccountIdPath + "/change-password"
)

const (
	// APIHealth represents the API to get the status of this application.
	APIHealth = API + "/health"
)
