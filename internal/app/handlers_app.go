package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
	"github.com/ptsypyshev/shortlink/internal/models"
)

const SessionAdmin = "admin"

func (a App) HandlerIndex(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	c.HTML(http.StatusOK, "main", gin.H{
		"title":         "Shortlink - make your links as short as possible!",
		"h1_text":       "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"page_template": "main",
		"is_admin":      userSession == SessionAdmin,
	})
}

func (a App) HandlerDashboard(c *gin.Context) {
	userID, session, userSession, err := a.HandlerDefault(c)
	if err != nil {
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
		"title":         "Shortlink - Dashboard for user %s",
		"h1_text":       "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"userID":        userID,
		"is_admin":      userSession == SessionAdmin,
		"page_template": "dashboard",
	})
}

func (a App) HandlerDefault(c *gin.Context) (any, sessions.Session, any, error) {
	userID, ok := c.Get("userID")
	if !ok {
		fmt.Println("cannot get userid")
	}
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession == nil {
		log.Println("Invalid session token")
		c.Redirect(http.StatusFound, "/login")
		return userID, session, userSession, fmt.Errorf("invalid session token")
	}
	return userID, session, userSession, nil
}

func (a App) HandlerUsersManagement(c *gin.Context) {
	userID, session, userSession, err := a.HandlerDefault(c)
	if err != nil {
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

	c.HTML(http.StatusOK, "users", gin.H{
		"title":         "Shortlink - User Management",
		"h1_text":       "Shortlink - make your links as short as possible!",
		"user_session":  userSession,
		"userID":        userID,
		"is_admin":      userSession == SessionAdmin,
		"page_template": "users",
	})
}

func (a App) HandlerLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession != nil {
		msg := "Please logout first"
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":         fmt.Sprintf("Shortlink - %s", msg),
				"user_session":  userSession,
				"page_template": "login",
				"error_message": msg,
				"do_logout":     "true",
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
		msg := "Please Sign out first"
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":         fmt.Sprintf("Shortlink - %s", msg),
				"user_session":  userSession,
				"page_template": "login",
				"error_message": msg,
				"do_logout":     "true",
			})
		return
	}

	user := &models.User{}
	if err := c.Bind(&user); err != nil {
		c.String(http.StatusBadRequest, "cannot bind to user: %s", err)
		return
	}
	checkedUser, ok := a.users.Check(a.ctx, user)
	if !ok {
		msg := "Bad login/password"
		c.HTML(http.StatusUnauthorized, "login",
			gin.H{
				"title":         fmt.Sprintf("Shortlink - %s", msg),
				"user_session":  userSession,
				"page_template": "login",
				"error_message": msg,
			})
		return
	}
	session.Set(UserKey, checkedUser.Username)
	if err := session.Save(); err != nil {
		msg := "Failed to save session"
		c.HTML(http.StatusInternalServerError, "login",
			gin.H{
				"title":         fmt.Sprintf("Shortlink - %s", msg),
				"user_session":  userSession,
				"page_template": "login",
				"error_message": msg,
			})
		return
	}
	c.Set("userID", checkedUser.ID)
	c.Redirect(http.StatusFound, "/dashboard")
}

func (a App) HandlerLogout(c *gin.Context) {
	session := sessions.Default(c)
	userSession := session.Get(UserKey)
	if userSession == nil {
		msg := "Please Sign in first"
		c.HTML(http.StatusBadRequest, "login",
			gin.H{
				"title":         fmt.Sprintf("Shortlink - %s", msg),
				"user_session":  userSession,
				"page_template": "login",
				"error_message": msg,
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
		"is_admin":      userSession == SessionAdmin,
		"page_template": "api",
	})
}

func (a App) HandlerNoRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error404", gin.H{
		"title":   "Shortlink - Page is not found",
		"h1_text": "Shortlink - make your links as short as possible!",
	})
}
