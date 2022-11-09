package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func GetRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("sessions", getTemplatePath("layout.html.tpl"), getTemplatePath("sessions.html.tpl"))
	r.AddFromFiles("lambdas", getTemplatePath("layout.html.tpl"), getTemplatePath("lambdas.html.tpl"))
	return r
}

func GetHTMLHandler(templateName string, templateData gin.H) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, templateName, templateData)
	}
}

func getTemplatePath(templateName string) string {
	return filepath.Join(os.Getenv("HTML_TEMPLATE_DIR"), templateName)
}
