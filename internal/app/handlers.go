package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
	"github.com/ptsypyshev/shortlink/internal/models"
)

const DefaultUserID = 1

func (a App) HandlerIndex(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	c.HTML(http.StatusOK, "main", gin.H{
		"title":         "Shortlink - make your links as short as possible!",
		"h1_text":       "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"page_template": "main",
	})
}

func (a App) HandlerDashboard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		fmt.Println("cannot get userid")
	}
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession == nil {
		log.Println("Invalid session token")
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user, err := a.users.Search(a.ctx, "username", userSession)
	if err != nil {
		fmt.Printf("cannot find user with name %s\n", userSession)
		session.Delete(UserKey)
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if userID == nil {
		userID = user[0].ID
	}

	c.HTML(http.StatusOK, "dashboard", gin.H{
		"title":        "Shortlink - Dashboard for user %s",
		"h1_text":      "Shortlink - make your links as short as possible!",
		"user_session": userSession,
		//"links":        links,
		"userID":        userID,
		"is_admin":      userSession == "admin",
		"page_template": "dashboard",
	})
}

func (a App) HandlerLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession != nil {
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":         "Shortlink - Please logout first",
				"user_session":  userSession,
				"page_template": "login",
			})
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"title":         "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"page_template": "login",
	})
}

func (a App) HandlerLogin(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession != nil {
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":         "Shortlink - Please logout first",
				"user_session":  userSession,
				"page_template": "login",
			})
		return
	}

	user := &models.User{}
	if err := c.Bind(&user); err != nil {
		c.String(http.StatusInternalServerError, "cannot bind to user: %s", err)
		return
	}
	checkedUser, ok := a.users.Check(a.ctx, user)
	if !ok {
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":        "Shortlink - Bad login/password",
				"user_session": userSession,
			})
		return
	}

	session.Set(UserKey, checkedUser.Username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login",
			gin.H{"title": "Failed to save session"})
		return
	}

	c.Set("userID", checkedUser.ID)

	c.Redirect(http.StatusFound, "/dashboard")
}

func (a App) HandlerLogout(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession == nil {
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title": "Shortlink - Please login first",
			})
		return
	}
	session.Delete(UserKey)
	if err := session.Save(); err != nil {
		log.Println("Failed to save session:", err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (a App) HandlerInitSchema(c *gin.Context) {
	if err := pgdb.InitSchema(c, a.pool); err != nil {
		a.logger.Error(fmt.Sprintf(`cannot init schema: %s`, err))
		c.JSON(http.StatusInternalServerError, gin.H{"result": "DB is not initialized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "DB Initialized"})
}

func (a App) HandlerAddDemoData(c *gin.Context) {
	if err := pgdb.AddDemoData(c, a.pool); err != nil {
		a.logger.Error(fmt.Sprintf(`cannot add demo data: %s`, err))
		c.JSON(http.StatusInternalServerError, gin.H{"result": "Demo data is not added"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "Demo data is added"})
}

func (a App) HandlerAPIHelp(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	c.HTML(http.StatusOK, "api", gin.H{
		"title":         "Shortlink - API Help",
		"h1_text":       "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"page_template": "api",
	})
}

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

func (a App) HandlerNoRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error404", gin.H{
		"title":   "Shortlink - Page is not found",
		"h1_text": "Shortlink - make your links as short as possible!",
	})
}
