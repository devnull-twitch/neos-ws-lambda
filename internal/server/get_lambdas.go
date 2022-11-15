package server

import (
	"net/http"
	"strings"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLambdas(c *gin.Context) {
	se, err := lambda.Get(c.Param("namespace"))
	if err != nil {
		logrus.WithError(err).Error("unable to load lambda session")
		c.Status(http.StatusInternalServerError)
		return
	}

	respStr := strings.Join(se.GetLambdaNames(), "\n")
	c.Data(http.StatusOK, "text/plain", []byte(respStr))
}
