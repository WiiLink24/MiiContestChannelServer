package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		authenticated := session.Get("authenticated")
		if authenticated == nil || !authenticated.(bool) {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/login")
			c.Abort()
			return
		}

		idToken := session.Get("id_token")
		if idToken == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
