package app

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserKey = "user"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user == nil {
		log.Println("User not logged in")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}
