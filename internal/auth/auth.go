package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.tor.ph/hiveon/pool/config"
	"git.tor.ph/hiveon/pool/internal/platform/redis"

	"github.com/sirupsen/logrus"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	codeRedirect2 = 302
	keyToken      = "accesstoken"
	keyNextPage   = "next"
)

var (
	// PathLogin is the path to handle OAuth 2.0 logins.
	PathLogin = "/login"
	// PathLogout is the path to handle OAuth 2.0 logouts.
	PathLogout = "/logout"
	// PathCallback is the path to handle callback from OAuth 2.0 backend
	// to exchange credentials.
	PathCallback = "/callback"
	// PathAdmin is the path to admin page
	PathAdmin = "/admin/"
	// PathError is the path to handle error cases.
	PathError         = "/unauthorized"
	oauthStateString2 = "state-string"
	conf              *oauth2.Config
	redirectURI       = config.Admin.HydraClient.CallbackURL
	log               *logrus.Logger
	e                 *casbin.Enforcer
)

var (
	ErrOAuthState    = errors.New("invalid oauth state")
	ErrTokenExchange = errors.New("token exchange failed")
	ErrUserInfo      = errors.New("cannot retrieve user info")
	ErrTokenInvalid  = errors.New("token invalid")
)

type User struct {
	ID   string `json:"sub"`
	Name string `json:"name"`
}

// Tokens represents a container that contains user's OAuth 2.0 access and refresh tokens.
type Tokens interface {
	Access() string
	Refresh() string
	Expired() bool
	ExpiryTime() time.Time
}

type token struct {
	oauth2.Token
}

// Access returns the access token.
func (t *token) Access() string {
	return t.AccessToken
}

// Refresh returns the refresh token.
func (t *token) Refresh() string {
	return t.RefreshToken
}

// Expired returns whether the access token is expired or not.
func (t *token) Expired() bool {
	if t == nil {
		return true
	}
	return !t.Token.Valid()
}

// ExpiryTime returns the expiry time of the user's access token.
func (t *token) ExpiryTime() time.Time {
	return t.Expiry
}

// String returns the string representation of the token.
func (t *token) String() string {
	return fmt.Sprintf("tokens: %s expire at: %s", t.Access(), t.ExpiryTime())
}

// IDP returns a new Google OAuth 2.0 backend endpoint.
func IDP(conf *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			switch c.Request.URL.Path {
			case PathLogin:
				login(conf, c)
			case PathLogout:
				logout(c)
			case PathCallback:
				handleOAuth2Callback(conf, c)
			}
		}
		s := sessions.Default(c)
		tk := unmarshallToken(s)

		if tk == nil {
			http.Redirect(c.Writer, c.Request, PathLogin, codeRedirect2)
			return
		}

		if tk != nil {
			// check if the access token is expired
			if tk.Expired() && tk.Refresh() == "" {
				s.Delete(keyToken)
				s.Save()
				tk = nil
			}
		}
		checkUser(c)
	}
}

// Handler that redirects user to the login page
// if user is not logged in.
// Sample usage:
// m.Get("/login-required", oauth2.LoginRequired, func() ... {})
var LoginRequired = func() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		url := conf.AuthCodeURL(oauthStateString2)
		token := unmarshallToken(s)
		if token == nil || token.Expired() {
			http.Redirect(c.Writer, c.Request, url, codeRedirect2)
		}
	}
}()

func login(f *oauth2.Config, c *gin.Context) {
	s := sessions.Default(c)
	url := conf.AuthCodeURL(oauthStateString2)
	if s.Get(keyToken) == nil {
		// User is not logged in.
		if url == "" {
			log.Error("can't redirect to empty url!")
			return
		}
		http.Redirect(c.Writer, c.Request, f.AuthCodeURL(url), codeRedirect2)
		return
	}
	// No need to login, redirect to the next page.
	http.Redirect(c.Writer, c.Request, url, codeRedirect2)
}

func logout(c *gin.Context) {
	s := sessions.Default(c)
	next := extractPath(c.Request.URL.Query().Get(keyNextPage))
	s.Delete(keyToken)
	s.Save()
	http.Redirect(c.Writer, c.Request, next, codeRedirect2)
}

func handleOAuth2Callback(f *oauth2.Config, c *gin.Context) {
	s := sessions.Default(c)
	code := c.Request.URL.Query().Get("code")
	t, err := f.Exchange(oauth2.NoContext, code)
	if err != nil {
		// Pass the error message, or allow dev to provide its own
		// error handler.
		log.Println("exchange oauth token failed:", err)
		http.Redirect(c.Writer, c.Request, PathError, codeRedirect2)
		return
	}
	// Store the credentials in the session.
	val, _ := json.Marshal(t)
	s.Set(keyToken, val)
	s.Save()
	http.Redirect(c.Writer, c.Request, PathAdmin, codeRedirect2)
}

func unmarshallToken(s sessions.Session) (t *token) {
	if s.Get(keyToken) == nil {
		return
	}
	data := s.Get(keyToken).([]byte)
	var tk oauth2.Token
	json.Unmarshal(data, &tk)
	return &token{tk}
}

func extractPath(next string) string {
	n, err := url.Parse(next)
	if err != nil {
		return "/"
	}
	return n.Path
}

func AuthConfig() *oauth2.Config {

	conf = &oauth2.Config{
		ClientID:     config.Admin.HydraClient.ClientID,
		ClientSecret: config.Admin.HydraClient.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Hydra.APIUrl + "/oauth2/auth",
			TokenURL: config.Hydra.APIUrl + "/oauth2/token",
		},
		RedirectURL: redirectURI,
		Scopes:      getScopes(),
	}
	return conf
}

func CheckAccess() gin.HandlerFunc {
	// casbin
	logrus.Info("created adapter and enforcer")
	a := redis.Adapter(config.Redis)
	e = NewAccessEnforcer(a)

	return IDP(AuthConfig())
}

func getScopes() []string {
	return []string{"openid", "offline"}
}

func GetCurrentUser(c *gin.Context) (string, error) {
	s := sessions.Default(c)

	if s.Get(keyToken) == nil {
		err := fmt.Errorf("invalid access token")
		log.Println("can't get current user: %s", err.Error())
		return "", err
	}

	token, errNoToken := getToken(c)

	if errNoToken != nil {
		log.Println("no token: %s", errNoToken)
		return "", errNoToken
	}

	httpClient := conf.Client(oauth2.NoContext, token)
	resp, err := httpClient.Get(config.Hydra.APIUrl + "/userinfo")

	if err != nil {
		log.Println("cannot get /userinfo: %s", err.Error())
		return "", err
	}

	defer resp.Body.Close()

	var user User

	errCantDecode := json.NewDecoder(resp.Body).Decode(&user)

	if errCantDecode != nil {
		return "", errCantDecode
	}

	return user.ID, nil
}

func getToken(c *gin.Context) (*oauth2.Token, error) {
	s := sessions.Default(c)

	accessToken := unmarshallToken(s)
	if accessToken == nil {
		return nil, ErrTokenInvalid
	}

	return &accessToken.Token, nil
}

func checkUser(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		http.Redirect(c.Writer, c.Request, PathLogin, codeRedirect2)
	}

	path := c.Request.URL.Path
	if strings.Contains(path, "admin") {
		e.LoadPolicy()
		if !e.Enforce(user, path, c.Request.Method) {
			http.Error(c.Writer, http.StatusText(403), 403)
			c.AbortWithStatus(403)
			return
		}
	}
}
