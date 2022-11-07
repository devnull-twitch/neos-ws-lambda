package server

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
)

func PostNamespace(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	args, err := parseArguments(string(bodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	se := lambda.NewEntry()
	token := token(10)
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
