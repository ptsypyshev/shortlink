package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
	"github.com/ptsypyshev/shortlink/internal/models"
	"github.com/ptsypyshev/shortlink/internal/repositories/objrepo"

	//nice "github.com/ekyoung/gin-nice-recovery"
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
	a.router = gin.Default()
	a.router.Static("/static", "./web/static")
	a.router.LoadHTMLGlob("web/templates/*")

	//Routes
	a.router.GET("/", a.HandlerIndex)
	a.router.GET("/:token", a.HandlerShortLink)
	a.router.GET("/login", a.HandlerLoginPage)
	a.router.POST("/login", a.HandlerLogin)
	a.router.GET("/api/", a.HandlerAPIHelp)
	a.router.GET("/dbinit/", a.HandlerInitSchema)
	a.router.GET("/demodb/", a.HandlerAddDemoData)

	a.router.GET("/api/users/:id", a.GetUser)
	a.router.POST("/api//users/", a.CreateUser)
	a.router.PUT("/api/users/", a.UpdateUser)
	a.router.DELETE("/api/users/:id", a.DeleteUser)

	a.router.GET("/api//links/:id", a.GetLink)
	a.router.POST("/api//links/", a.CreateLink)
	a.router.PUT("/api//links/", a.UpdateLink)
	a.router.DELETE("/api//links/:id", a.DeleteLink)

	// Start serving the application
	return a.router.Run()
}
