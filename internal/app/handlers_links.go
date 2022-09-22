package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ptsypyshev/shortlink/internal/models"
)

const DefaultUserID = 1

func (a App) CreateLink(c *gin.Context) {
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if link.OwnerID == 0 {
		session := sessions.Default(c)
		userSession := session.Get(UserKey)
		if userSession == nil {
			link.OwnerID = DefaultUserID
		} else {
			user, err := a.users.Search(a.ctx, "username", userSession)
			if err != nil {
				fmt.Printf("cannot find user with name %s\n", userSession)
				link.OwnerID = DefaultUserID
			}
			link.OwnerID = user[0].ID
		}
	}

	newLink, err := a.links.Create(a.ctx, &link)
	if err != nil {
		msg := fmt.Sprintf(`create link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	longLinkID := newLink.ID
	shortlink, err := a.shortlinks.Create(a.ctx, longLinkID)
	if err != nil {
		msg := fmt.Sprintf(`create shortlink error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": shortlink.Token})
}

func (a App) GetLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	link, err := a.links.Read(a.ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"read": link})
}

func (a App) UpdateLink(c *gin.Context) {
	var link models.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	id := link.ID
	updatedLink, err := a.links.Update(a.ctx, id, &link)
	if err != nil {
		msg := fmt.Sprintf(`update link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": updatedLink})
}

func (a App) DeleteLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	deletedLink, err := a.links.Delete(a.ctx, id)
	if err != nil {
		msg := fmt.Sprintf(`delete link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deletedLink})
}

func (a App) SearchLinks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf(`bad id: %s`, c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	foundLinks, err := a.links.Search(a.ctx, "owner_id", id)
	if err != nil {
		msg := fmt.Sprintf(`no links found: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"found": foundLinks})
}

func (a App) HandlerShortLink(c *gin.Context) {
	value := c.Param("token")
	shortlinks, err := a.shortlinks.Search(a.ctx, "token", value)
	if err != nil {
		msg := fmt.Sprintf(`get error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	if len(shortlinks) != 1 {
		a.HandlerNoRoute(c)
		return
	}
	linkID := shortlinks[0].LongLinkID
	link, err := a.links.Read(a.ctx, linkID)
	if err != nil {
		msg := fmt.Sprintf(`read link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	link.ClickCounter++
	link, err = a.links.Update(a.ctx, linkID, link)
	if err != nil {
		msg := fmt.Sprintf(`update link error: %s`, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.Redirect(http.StatusFound, link.LongLink)
}
