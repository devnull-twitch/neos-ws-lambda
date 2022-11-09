package server

import (
	"fmt"
	"net/http"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PostSaveSession(c *gin.Context) {
	se, err := lambda.Get(c.Param("namespace"))
	if err != nil {
		logrus.WithError(err).Error("unable to load session")
		c.Status(http.StatusInternalServerError)
		return
	}

	tpl := se.ToTemplate()
	token := token(10)
	if err := lambda.WriteTemplate(fmt.Sprintf("%s.json", token), tpl); err != nil {
		logrus.WithError(err).Error("unable to save session as template")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusCreated, "text/plain", []byte(token))
}
