package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
	"github.com/ptsypyshev/shortlink/internal/models"
	"github.com/ptsypyshev/shortlink/internal/repositories/objrepo"

	//nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	ctx        context.Context
	router     *gin.Engine
	pool       *pgxpool.Pool
	users      objrepo.Users
	links      objrepo.Links
	shortlinks objrepo.ShortLinks
	logger     *zap.Logger
	//tracer   opentracing.Tracer
}

// func (a *App) Init() (io.Closer, error) {

func (a *App) Init() error {
	a.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot init Logger: %s", err)
	}
	defer func() { _ = logger.Sync() }()
	//tracer, closer := InitJaeger(logger)
	//tracer, closer := InitJaeger("Shortlink", "localhost:6831", logger)

	pool, err := pgdb.InitDB(a.ctx, logger)
	if err != nil {
		log.Fatalf("cannot init DB: %s", err)
	}

	if _, err := os.Stat("/configured"); errors.Is(err, os.ErrNotExist) {
		if err := pgdb.InitSchema(a.ctx, pool); err != nil {
			a.logger.Error(fmt.Sprintf(`cannot init schema: %s`, err))
			log.Fatalf("cannot init DB: %s", err)
		}
		file, err := os.Create("/configured")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	UsersDB := pgdb.DBNew[*models.User](pool)
	LinksDB := pgdb.DBNew[*models.Link](pool)
	ShortLinksDB := pgdb.DBNew[*models.ShortLink](pool)

	users := objrepo.UsersNew(UsersDB, logger)
	links := objrepo.LinksNew(LinksDB, logger)
	shortlinks := objrepo.ShortLinksNew(ShortLinksDB, logger)

	a.logger = logger
	a.pool = pool
	a.users = *users
	a.links = *links
	a.shortlinks = *shortlinks
	return nil
	//return closer, nil
}

func (a *App) Serve() error {
	//Initialize Router and add Middleware
	a.router = gin.New()
	a.router.Static("/static", "./web/static")
	a.router.LoadHTMLGlob("web/templates/*")
	a.router.Use(sessions.Sessions("session", cookie.NewStore([]byte("secret"))))

	//Routes
	public := a.router.Group("/")
	{
		public.GET("/", a.HandlerIndex)
		public.GET("/:token", a.HandlerShortLink)
		//public.GET("/api/", a.HandlerAPIHelp)
		public.GET("/login", a.HandlerLoginPage)
		public.POST("/login", a.HandlerLogin)
		public.POST("/api/links/", a.CreateLink)
	}

	private := a.router.Group("/")
	private.Use(AuthRequired)
	{
		private.GET("/dbinit/", a.HandlerInitSchema)
		private.GET("/demodb/", a.HandlerAddDemoData)
		private.GET("/api/", a.HandlerAPIHelp)
		private.GET("/dashboard/", a.HandlerDashboard)
		private.GET("/logout/", a.HandlerLogout)

		private.GET("/api/users/:id", a.GetUser)
		private.POST("/api//users/", a.CreateUser)
		private.PUT("/api/users/", a.UpdateUser)
		private.DELETE("/api/users/:id", a.DeleteUser)
		private.GET("/api/users/:id/links", a.SearchLinks)

		private.GET("/api/links/:id", a.GetLink)
		//private.POST("/api//links/", a.CreateLink)
		private.PUT("/api/links/", a.UpdateLink)
		private.DELETE("/api/links/:id", a.DeleteLink)
	}

	return a.router.Run()
}
