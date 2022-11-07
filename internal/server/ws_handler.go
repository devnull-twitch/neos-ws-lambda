package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

type LambdaArgs map[string]string

func WsHandler(c *gin.Context) {
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

		msg := string(rb)
		if msg == "" {
			logrus.Warn("empty message from ws")
			continue
		}

		fnName, arguments, err := parseMessage(msg)
		if err != nil {
			logrus.WithError(err).Error("empty message from ws")
			continue
		}

		se.RunLambda(fnName, arguments)
		namespaceLog.WithField("msg", msg).Info("exec lambda")
	}
}

func parseMessage(msg string) (string, LambdaArgs, error) {
	firstSepIndex := strings.Index(msg, "|")
	lambdaName := msg

	var arguments map[string]string
	var err error
	if firstSepIndex > 0 {
		lambdaName = msg[0:firstSepIndex]

		arguments, err = parseArguments(msg[firstSepIndex+1:])
		if err != nil {
			return "", nil, err
		}
	}

	return lambdaName, arguments, nil
}

func parseArguments(msg string) (LambdaArgs, error) {
	parts := strings.Split(msg, "|")
	arguments := make(map[string]string)
	for _, part := range parts {
		firstEqualSignIndex := strings.Index(part, "=")
		if firstEqualSignIndex == -1 {
			return nil, fmt.Errorf("missing equal sign in argument")
		}

		arguments[part[0:firstEqualSignIndex]] = part[firstEqualSignIndex+1:]
	}

	return arguments, nil
}
