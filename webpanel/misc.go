package webpanel

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (w *WebPanel) MarqueePage(c *gin.Context) {
	c.HTML(http.StatusOK, "marquee.html", gin.H{})
}

func (w *WebPanel) EditMarquee(c *gin.Context) {
	marquee_text := c.PostForm("marquee_text")

	err := os.WriteFile(fmt.Sprintf("%s/marquee/marquee.txt", w.Config.AssetsPath), []byte(marquee_text), 0666)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/panel/marquee")
}
