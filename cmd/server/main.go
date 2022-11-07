package main

import (
	"time"

	"github.com/devnull-twitch/neos-ws-lambda/internal/server"
	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/lambda", server.PostNamespace)
	r.POST("/lambda/:namespace/:func", server.PostLambda)
	r.GET("/connect/:namespace", server.WsHandler)

	doneChan := make(chan bool)
	go lambda.CleanupWorker(doneChan)

	r.Run(":8081")

	select {
	case doneChan <- true:
	case <-time.After(time.Second):
	}
}
