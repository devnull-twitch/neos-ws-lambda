package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostSession(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	se := lambda.NewEntry()
	token := token(10)

	if c.Query("template") != "" {
		tplName := c.Query("template")
		validate := validator.New()
		if err := validate.Var(tplName, "alphanum"); err != nil {
			log.Fatal(err)
		}
		tpl, err := lambda.ReadTemplate(fmt.Sprintf("%s.json", tplName))
		if err != nil {
			log.Fatal(err)
		}

		for tplVarName, tplVarVal := range tpl.Arguments {
			se.SetupPersist(tplVarName, tplVarVal)
		}

		for tplLambdaName, tplLambdaCode := range tpl.Lambdas {
			se.AddLambda(tplLambdaName, tplLambdaCode)
		}
	}

	args, err := parseArguments(string(bodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	for varName, varVal := range args {
		se.SetupPersist(varName, varVal)
	}

	lambda.Add(token, se)

	c.Data(http.StatusCreated, "text/plain", []byte(token))
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func tokenWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func token(length int) string {
	return tokenWithCharset(length, charset)
}
