package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	oauth_gh "golang.org/x/oauth2/github"
)

var signKey = []byte(os.Getenv("JWT_SECRET"))

func GetGHReditectHandler(addStateChan chan<- string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := getConf()
		state := token(10)
		select {
		case addStateChan <- state:
			c.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(state))
		case <-time.After(time.Second):
			logrus.Error("github state add timeout")
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

}

type claims struct {
	jwt.RegisteredClaims
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

func GetGHCallbackHandler(stateCheckChan chan<- checkStateRequest) gin.HandlerFunc {
	conf := getConf()
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		resp := make(chan bool)
		select {
		case stateCheckChan <- checkStateRequest{
			state: state,
			resp:  resp,
		}:
		case <-time.After(time.Second):
			logrus.Error("github state check request timeout")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		select {
		case ok := <-resp:
			if !ok {
				c.AbortWithStatus(http.StatusBadRequest)
				logrus.Error("invalid state from github auth response")
				return
			}
		case <-time.After(time.Second):
			logrus.Error("github state check response timeout")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			logrus.WithError(err).Error("unable to exchange gh code for full token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		clientContext := context.Background()
		tSrc := conf.TokenSource(clientContext, tok)
		httpClient := oauth2.NewClient(clientContext, tSrc)
		ghClient := github.NewClient(httpClient)

		authedUser, _, err := ghClient.Users.Get(context.Background(), "")
		if err != nil {
			logrus.WithError(err).Error("unable to load authed user from api after getting token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tokenClaims := claims{
			jwt.RegisteredClaims{
				Issuer:    "llws",
				Subject:   authedUser.GetEmail(),
				ExpiresAt: &jwt.NumericDate{Time: tok.Expiry},
				IssuedAt:  &jwt.NumericDate{Time: time.Now()},
				ID:        fmt.Sprintf("%v", authedUser.GetID()),
			},
			tok.AccessToken,
			tok.TokenType,
			tok.RefreshToken,
			tok.Expiry,
		}

		fullJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
		signedToken, err := fullJwt.SignedString(signKey)
		if err != nil {
			logrus.WithError(err).Error("unable to load authed user from api after getting token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt", signedToken, int(tok.Expiry.Sub(time.Now())), "/", "", false, true)
		c.Status(http.StatusOK)
	}
}

func getConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Scopes:       []string{},
		Endpoint:     oauth_gh.Endpoint,
		RedirectURL:  os.Getenv("GH_CLIENT_REDIRECT"),
	}
}

var states = make([]string, 0)

type checkStateRequest struct {
	state string
	resp  chan<- bool
}

func GhStateWorker() (chan<- string, chan<- checkStateRequest) {
	addChan := make(chan string)
	checkChan := make(chan checkStateRequest)

	go func() {
		for {
			select {
			case newState := <-addChan:
				states = append(states, newState)
			case req := <-checkChan:
				filtered := make([]string, 0, len(states))
				found := false
				for _, savedState := range states {
					if savedState == req.state {
						found = true
					} else {
						filtered = append(filtered, savedState)
					}
				}
				if found {
					states = filtered
				}
				req.resp <- found
			}
		}
	}()

	return addChan, checkChan
}
