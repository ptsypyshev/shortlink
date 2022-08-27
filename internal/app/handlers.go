package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
	"github.com/ptsypyshev/shortlink/internal/models"
)

func (a App) HandlerIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "main", gin.H{
		"title":   "Shortlink - make your links as short as possible!",
		"h1_text": "Shortlink - make your links as short as possible!",
	})
}

func (a App) HandlerLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Shortlink - make your links as short as possible!",
	})
}

func (a App) HandlerLogin(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(&user); err != nil {
		c.String(http.StatusInternalServerError, "cannot bind to user: %s", err)
		return
	}
	c.String(http.StatusOK, "User is %v\n", user)
	//c.HTML(http.StatusOK, "login", gin.H{
	//	"title": "Shortlink - make your links as short as possible!",
	//})
}

func (a App) HandlerInitSchema(c *gin.Context) {
	if err := pgdb.InitSchema(c, a.pool); err != nil {
		a.logger.Error(fmt.Sprintf(`cannot init schema: %s`, err))
		c.String(http.StatusInternalServerError, "DB is not initialized")
		return
	}
	c.String(http.StatusOK, "DB Initialized")
}

func (a App) HandlerAddDemoData(c *gin.Context) {
	if err := pgdb.AddDemoData(c, a.pool); err != nil {
		a.logger.Error(fmt.Sprintf(`cannot add demo data: %s`, err))
		c.String(http.StatusInternalServerError, "Demo data is not added")
		return
	}
	c.String(http.StatusOK, "Demo data is added")
}

func (a App) HandlerAPIHelp(c *gin.Context) {
	c.HTML(http.StatusOK, "api", gin.H{
		"title":   "Shortlink - API Help",
		"h1_text": "Shortlink - make your links as short as possible!",
	})
}

func (a App) CreateUser(c *gin.Context) {
	ctx := context.Background()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	newUser, err := a.users.Create(ctx, &user)
	if err != nil {
		msg := fmt.Sprintf(`create user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": newUser})
}

func (a App) GetUser(c *gin.Context) {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	user, err := a.users.Read(ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"read": user})
}

func (a App) UpdateUser(c *gin.Context) {
	ctx := context.Background()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	id := user.ID
	updatedUser, err := a.users.Update(ctx, id, &user)
	if err != nil {
		msg := fmt.Sprintf(`update user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": updatedUser})
}

func (a App) DeleteUser(c *gin.Context) {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	deletedUser, err := a.users.Delete(ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`delete user error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deletedUser})
}

func (a App) CreateLink(c *gin.Context) {
	ctx := context.Background()
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if link.OwnerID == 0 {
		link.OwnerID = 1
	}
	newLink, err := a.links.Create(ctx, &link)
	if err != nil {
		msg := fmt.Sprintf(`create link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	longLinkID := newLink.ID
	shortlink, err := a.shortlinks.Create(ctx, longLinkID)
	if err != nil {
		msg := fmt.Sprintf(`create shortlink error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": shortlink.Token})
}

func (a App) GetLink(c *gin.Context) {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	link, err := a.links.Read(ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"read": link})
}

func (a App) UpdateLink(c *gin.Context) {
	ctx := context.Background()
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	id := link.ID
	updatedLink, err := a.links.Update(ctx, id, &link)
	if err != nil {
		msg := fmt.Sprintf(`update link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": updatedLink})
}

func (a App) DeleteLink(c *gin.Context) {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	deletedLink, err := a.links.Delete(ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`delete link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deletedLink})
}

func (a App) HandlerShortLink(c *gin.Context) {
	ctx := context.Background()
	value := c.Param("token")
	shortlinks, err := a.shortlinks.Search(ctx, "token", value)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	if len(shortlinks) != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get bad redirect for %s", value)})
		return
	}
	linkID := shortlinks[0].LongLinkID
	link, err := a.links.Read(ctx, linkID)
	if err != nil {
		msg := fmt.Sprintf(`read link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	link.ClickCounter++
	link, err = a.links.Update(ctx, linkID, link)
	if err != nil {
		msg := fmt.Sprintf(`update link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.Redirect(http.StatusFound, link.LongLink)
}
