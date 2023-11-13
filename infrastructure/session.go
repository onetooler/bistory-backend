package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/logger"
	"gopkg.in/boj/redistore.v1"
)

// TODO: change gorilla/sessions to SCS for better performance

const (
	// sessionStr represents a string of session key.
	sessionStr = "GSESSION"
	// accountStr is the key of account data in the session.
	accountStr = "Account"
	// emailVerification is the key of email verification token in the session.
	emailVerificationStr = "EmailVerification"
)

type session struct {
	store sessions.Store
}

// Session represents a interface for accessing the session on the application.
type Session interface {
	GetStore() sessions.Store

	Get(c echo.Context) *sessions.Session
	Save(c echo.Context) error
	Delete(c echo.Context) error
	SetValue(c echo.Context, key string, value string) error
	GetValue(c echo.Context, key string) string
	SetAccount(c echo.Context, account *Account) error
	GetAccount(c echo.Context) *Account
	SetEmailVerification(c echo.Context, emailVerification *EmailVerification) error
	VerifyEmailToken(c echo.Context, token string) error
	IsVerifiedEmail(c echo.Context, email string) (bool, error)
	Login(c echo.Context, account *Account) error
	Logout(c echo.Context) error
	HasAuthorizationTo(c echo.Context, accountId uint, authority uint) bool
}

type Account struct {
	Id        uint      `json:"id"`
	LoginId   string    `json:"loginId"`
	LoginTime time.Time `json:"loginTime"`
	Authority uint      `json:"authority"`
}

type EmailVerification struct {
	Email            string    `json:"email"`
	Token            string    `json:"token"`
	TokenGeneratedAt time.Time `json:"tokenGeneratedAt"`
	VerifiedAt       time.Time `json:"verifiedAt"`
}

// NewSession is constructor.
func NewSession(logger logger.Logger, conf *config.Config) Session {
	if !conf.Redis.Enabled {
		logger.GetZapLogger().Infof("use CookieStore for session")
		return &session{sessions.NewCookieStore([]byte("secret"))}
	}

	logger.GetZapLogger().Infof("use redis for session")
	logger.GetZapLogger().Infof("Try redis connection")
	address := fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port)
	store, err := redistore.NewRediStore(conf.Redis.ConnectionPoolSize, "tcp", address, "", []byte("secret"))
	if err != nil {
		logger.GetZapLogger().Panicf("Failure redis connection, %s", err.Error())
	}
	logger.GetZapLogger().Infof(fmt.Sprintf("Success redis connection, %s", address))
	return &session{store: store}
}

func (s *session) GetStore() sessions.Store {
	return s.store
}

// Get returns a session for the current request.
func (s *session) Get(c echo.Context) *sessions.Session {
	sess, _ := s.store.Get(c.Request(), sessionStr)
	return sess
}

// Save saves the current session.
func (s *session) Save(c echo.Context) error {
	sess := s.Get(c)
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	return s.saveSession(c, sess)
}

// Delete the current session.
func (s *session) Delete(c echo.Context) error {
	sess := s.Get(c)
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	return s.saveSession(c, sess)
}

func (s *session) saveSession(c echo.Context, sess *sessions.Session) error {
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error occurred while save session")
	}
	return nil
}

// SetValue sets a key and a value.
func (s *session) SetValue(c echo.Context, key string, value string) error {
	sess := s.Get(c)
	sess.Values[key] = value
	return nil
}

// GetValue returns value of session.
func (s *session) GetValue(c echo.Context, key string) string {
	sess := s.Get(c)
	if sess != nil {
		if v, ok := sess.Values[key]; ok {
			data, result := v.(string)
			if result && data != "null" {
				return data
			}
		}
	}
	return ""
}

func (s *session) Login(c echo.Context, account *Account) error {
	account.LoginTime = time.Now()
	if err := s.SetAccount(c, account); err != nil {
		return err
	}
	if err := s.Save(c); err != nil {
		return err
	}
	return nil
}

func (s *session) Logout(c echo.Context) error {
	if err := s.SetAccount(c, nil); err != nil {
		return err
	}
	if err := s.Delete(c); err != nil {
		return err
	}
	return nil
}

func (s *session) SetAccount(c echo.Context, account *Account) error {
	bytes, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("json marshal error while set value in session")
	}
	return s.SetValue(c, accountStr, string(bytes))
}

func (s *session) GetAccount(c echo.Context) *Account {
	if v := s.GetValue(c, accountStr); v != "" {
		a := &Account{}
		_ = json.Unmarshal([]byte(v), a)
		return a
	}
	return nil
}

func (s *session) SetEmailVerification(c echo.Context, emailVerification *EmailVerification) error {
	bytes, err := json.Marshal(emailVerification)
	if err != nil {
		return fmt.Errorf("json marshal error while set value in session")
	}

	if err := s.SetValue(c, emailVerificationStr, string(bytes)); err != nil {
		return err
	}
	return s.Save(c)
}

func (s *session) GetEmailVerification(c echo.Context) *EmailVerification {
	if v := s.GetValue(c, emailVerificationStr); v != "" {
		e := &EmailVerification{}
		_ = json.Unmarshal([]byte(v), e)
		return e
	}
	return nil
}

func (s *session) VerifyEmailToken(c echo.Context, token string) error {
	emailVerification := s.GetEmailVerification(c)
	if emailVerification == nil {
		return fmt.Errorf("emailVerification not found")
	}
	if emailVerification.TokenGeneratedAt.Before(time.Now().Add(-config.EmailVerificationTokenLifetime)) {
		_ = s.SetEmailVerification(c, nil)
		return fmt.Errorf("token expired")
	}
	if token != emailVerification.Token {
		return fmt.Errorf("token not matched")
	}
	emailVerification.VerifiedAt = time.Now()
	return s.SetEmailVerification(c, emailVerification)
}

func (s *session) IsVerifiedEmail(c echo.Context, email string) (bool, error) {
	emailVerification := s.GetEmailVerification(c)
	if emailVerification.VerifiedAt.IsZero() {
		return false, fmt.Errorf("not verified yet")
	}
	if emailVerification.VerifiedAt.Before(time.Now().Add(-config.EmailVerificationLifetime)) {
		_ = s.SetEmailVerification(c, nil)
		return false, fmt.Errorf("verification expired")
	}
	if email != emailVerification.Email {
		return false, fmt.Errorf("email not matched")
	}
	return true, nil
}

func (s *session) HasAuthorizationTo(c echo.Context, accountId uint, authority uint) bool {
	currentAccount := s.GetAccount(c)
	if currentAccount == nil {
		return false
	}
	if currentAccount.Id == accountId {
		return true
	}
	if currentAccount.Authority < authority {
		return true
	}
	return false
}
