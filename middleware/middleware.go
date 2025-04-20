package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/login")
			c.Abort()
			return
		}

		// Verify the OpenID Connect idToken.
		ctx := context.Background()
		idToken, err := verifier.Verify(ctx, tokenString)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/login")
			c.Abort()
			return
		}

		// Parse custom claims if needed.
		var claims struct {
			Username string `json:"username"`
			UserId   string `json:"user_id"`
		}
		if err = idToken.Claims(&claims); err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/login")
			c.Abort()
			return
		}
		log.Println(claims.Username, claims.UserId)

		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
