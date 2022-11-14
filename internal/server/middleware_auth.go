package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func CheckAuth(c *gin.Context) {
	jwtStr, err := c.Cookie("jwt")
	if err != nil {
		logrus.WithError(err).Warn("unable to read jwt")
		return
	}
	if len(jwtStr) <= 0 {
		return
	}

	token, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid jwt signature method")
		}

		return signKey, nil
	})
	if err != nil {
		logrus.WithError(err).Error("jwt parsing/validation error")
		c.SetCookie("jwt", "", -1, "/", "", false, true)
		return
	}

	claims := token.Claims.(claims)
	c.Set("login", claims)
}
