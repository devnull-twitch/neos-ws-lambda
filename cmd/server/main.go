package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := gin.Default()
	r.POST("/lambda/:namespace", func(c *gin.Context) {
		persistVars := make([]string, 0)
		if err := c.BindJSON(&persistVars); err != nil {
			log.Fatal(err)
		}

		se := lambda.NewEntry()
		se.Namespace = c.Param("namespace")
		for _, varName := range persistVars {
			se.SetupPersist(varName)
		}

		lambda.Add(c.Param("namespace"), se)
	})
	r.POST("/lambda/:namespace/:func", func(c *gin.Context) {
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err)
		}

		se, err := lambda.Get(c.Param("namespace"))
		if err != nil {
			log.Fatal(err)
		}

		se.AddLambda(c.Param("func"), string(bodyBytes))
	})
	r.GET("/connect/:namespace", func(c *gin.Context) {
		se, err := lambda.Get(c.Param("namespace"))
		if err != nil {
			log.Fatal(err)
		}

		namespaceLog := logrus.WithField("namespace", c.Param("namespace"))

		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("error get connection")
			log.Fatal(err)
		}
		defer ws.Close()

		writeChannel := make(chan lambda.MessageTpl)
		endChan := make(chan bool)
		go func(wc <-chan lambda.MessageTpl, done <-chan bool) {
			run := true
			logrus.Info("started ws writer")
			for run {
				select {
				case msg := <-wc:
					transmitData := fmt.Sprintf("%s|%v", msg.VarName, msg.VarVal)
					err := ws.WriteMessage(websocket.TextMessage, []byte(transmitData))
					if err != nil {
						log.Fatal(err)
					}
				case <-done:
					run = false
				}
			}
			logrus.Info("ended ws writer")
		}(writeChannel, endChan)
		se.SetWriterChannel(writeChannel)

		defer func() {
			select {
			case endChan <- true:
			case <-time.After(time.Second):
			}
		}()

		for {
			mt, rb, err := ws.ReadMessage()
			if err != nil {
				break
			}

			if mt != websocket.TextMessage {
				continue
			}

			se.RunLambda(string(rb))
			namespaceLog.WithField("lambda", string(rb)).Info("exec lambda")
		}
	})

	r.Run(":8081")
}
