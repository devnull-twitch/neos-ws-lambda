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
		// memory writing endpoints
		api.POST("/session", server.PostSession)
		api.POST("/session/:namespace/lambda/:func", server.PostLambda)
		api.GET("/session/:namespace/lambda", server.GetLambdas)
	}

	r.GET("/connect/:namespace", server.WsHandler)

	// HTML pages
	r.GET("/", server.GetHTMLHandler("sessions", gin.H{}))
	r.GET("/session/:namespace", func(c *gin.Context) {
		server.GetHTMLHandler("lambdas", gin.H{
			"Namespace": c.Param("namespace"),
		})(c)
	})

	r.Static("/assets", os.Getenv("ASSET_DIR"))

	ghAddStateChan, ghCheckChan := server.GhStateWorker()

	ghAuthGroup := r.Group("/auth/github")
	{
		ghAuthGroup.GET("/", server.GetGHReditectHandler(ghAddStateChan))
		ghAuthGroup.GET("/back", server.GetGHCallbackHandler(ghCheckChan))
	}

	doneChan := make(chan bool)
	go lambda.CleanupWorker(doneChan)

	r.Run(":8081")

	select {
	case doneChan <- true:
	case <-time.After(time.Second):
	}
}
