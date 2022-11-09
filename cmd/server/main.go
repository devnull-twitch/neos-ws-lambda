package main

import (
	"os"
	"time"

	"github.com/devnull-twitch/neos-ws-lambda/internal/server"
	"github.com/devnull-twitch/neos-ws-lambda/lib/lambda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := gin.Default()
	r.HTMLRender = server.GetRenderer()

	api := r.Group("/api")
	{
		api.POST("/lambda", server.PostNamespace)
		api.POST("/lambda/:namespace/:func", server.PostLambda)
	}

	r.GET("/connect/:namespace", server.WsHandler)

	// HTML pages
	r.GET("/", server.GetHTMLHandler("sessions", gin.H{}))
	r.GET("/lambdas/:namespace", func(c *gin.Context) {
		server.GetHTMLHandler("lambdas", gin.H{
			"Namespace": c.Param("namespace"),
		})(c)
	})

	r.Static("/assets", os.Getenv("ASSET_DIR"))

	doneChan := make(chan bool)
	go lambda.CleanupWorker(doneChan)

	r.Run(":8081")

	select {
	case doneChan <- true:
	case <-time.After(time.Second):
	}
}
