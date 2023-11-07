package infrastructure

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	echoSession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// TODO: change gorilla/sessions to SCS for better performance

const (
	// sessionStr represents a string of session key.
	sessionStr = "GSESSION"
	// accountStr is the key of account data in the session.
	accountStr = "Account"
)

type session struct {
	context echo.Context
}

// Session represents a interface for accessing the session on the application.
type Session interface {
	SetContext(c echo.Context)
	Get() *sessions.Session
	Save() error
	Delete() error
	SetValue(key string, value interface{}) error
	GetValue(key string) string
	SetAccount(account *Account) error
	GetAccount() *Account
	Login(account *Account) error
	Logout() error
	HasAuthorizationTo(uint, uint) bool
}

type Account struct {
	Id        uint      `json:"id"`
	LoginId   string    `json:"loginId"`
	LoginTime time.Time `json:"loginTime"`
	Authority uint      `json:"authority"`
}

// NewSession is constructor.
func NewSession() Session {
	return &session{context: nil}
}

// SetContext sets the context of echo framework to the session.
func (s *session) SetContext(c echo.Context) {
	s.context = c
}

// Get returns a session for the current request.
func (s *session) Get() *sessions.Session {
	sess, _ := echoSession.Get(sessionStr, s.context)
	return sess
}

// Save saves the current session.
func (s *session) Save() error {
	sess := s.Get()
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	return s.saveSession(sess)
}

// Delete the current session.
func (s *session) Delete() error {
	sess := s.Get()
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	return s.saveSession(sess)
}

func (s *session) saveSession(sess *sessions.Session) error {
	if err := sess.Save(s.context.Request(), s.context.Response()); err != nil {
		return s.context.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// SetValue sets a key and a value.
func (s *session) SetValue(key string, value interface{}) error {
	sess := s.Get()
	bytes, err := json.Marshal(value)
	if err != nil {
		return s.context.NoContent(http.StatusInternalServerError)
	}
	sess.Values[key] = string(bytes)
	return nil
}

// GetValue returns value of session.
func (s *session) GetValue(key string) string {
	sess := s.Get()
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

func (s *session) Login(account *Account) error {
	account.LoginTime = time.Now()
	if err := s.SetAccount(account); err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

func (s *session) Logout() error {
	if err := s.SetAccount(nil); err != nil {
		return err
	}
	if err := s.Delete(); err != nil {
		return err
	}
	return nil
}

func (s *session) SetAccount(account *Account) error {
	return s.SetValue(accountStr, account)
}

func (s *session) GetAccount() *Account {
	if v := s.GetValue(accountStr); v != "" {
		a := &Account{}
		_ = json.Unmarshal([]byte(v), a)
		return a
	}
	return nil
}

func (s *session) HasAuthorizationTo(accountId uint, authority uint) bool {
	currentAccount := s.GetAccount()
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
