package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
)

func PostTemplate(c *gin.Context) {
	tpl := &lambda.Template{}
	if err := c.BindJSON(tpl); err != nil {
		log.Fatal(err)
	}

	token := token(10)
	if err := lambda.WriteTemplate(fmt.Sprintf("%s.json", token), tpl); err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusCreated, "text/plain", []byte(token))
}
