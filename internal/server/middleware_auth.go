package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	basicUser, basicPw, hasBasic := c.Request.BasicAuth()
	if !hasBasic {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if basicUser != os.Getenv("AUTH_USERNAME") || basicPw != os.Getenv("AUTH_PASSWORD") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
