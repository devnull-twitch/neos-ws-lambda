package server

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PostLambda(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	se, err := lambda.Get(c.Param("namespace"))
	if err != nil {
		logrus.WithError(err).Error("unable to load lambda session")
		c.Status(http.StatusInternalServerError)
		return
	}

	se.AddLambda(c.Param("func"), string(bodyBytes))
}
