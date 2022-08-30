package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ptsypyshev/shortlink/internal/models"
)

func (a App) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	newUser, err := a.users.Create(a.ctx, &user)
	if err != nil {
		msg := fmt.Sprintf(`create user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": newUser})
}

func (a App) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	user, err := a.users.Read(a.ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"read": user})
}

func (a App) GetUsers(c *gin.Context) {
	users, err := a.users.Search(a.ctx, "all", "")
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"found": users})
}

func (a App) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	id := user.ID
	updatedUser, err := a.users.Update(a.ctx, id, &user)
	if err != nil {
		msg := fmt.Sprintf(`update user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": updatedUser})
}

func (a App) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	deletedUser, err := a.users.Delete(a.ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`delete user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deletedUser})
}
