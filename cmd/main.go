/*
 * Shortlink Service
 *
 * It is an URL shortener service - it can generate a short version of arbitrary user's URL.
 * Also user can get statistics for URL generated url (how much people followed the link).
 *
 * API version: 1.0.0
 * Contact: ptsypyshev@gmail.com
 */
package main

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/ptsypyshev/shortlink/internal/db/pg"
	"github.com/ptsypyshev/shortlink/internal/models"
	"github.com/ptsypyshev/shortlink/internal/repositories/objrepo"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot initialize zap logger")
	}
	defer func() { _ = logger.Sync() }()
	//tracer, closer := InitJaeger(logger)
	//tracer, closer := InitJaeger("Shortlink", "localhost:6831", logger)
	//defer closer()

	pool, err := pg.InitDB(ctx, logger)
	if err != nil {
		log.Fatalf("cannot init DB: %s", err)
	}

	pg.InitSchema(ctx, pool)
	pg.AddDemoData(ctx, pool)

	UsersDB := pg.PgDBNew[*models.User](pool)
	EmptyUser := &models.User{}
	readUser, _ := UsersDB.Read(ctx, 2, EmptyUser)
	fmt.Println(readUser.String())
	id, _ := UsersDB.Create(ctx, readUser)
	fmt.Printf("User with id %d is created.\n", id)
	readUser.Username = "Updated"
	updatedUser, _ := UsersDB.Update(ctx, readUser)
	fmt.Printf("User with id %d is updated.\n%v\n", updatedUser.ID, updatedUser)
	UsersDB.Delete(ctx, 5)

	LinksDB := pg.PgDBNew[*models.Link](pool)
	EmptyLink := &models.Link{}
	readLink, _ := LinksDB.Read(ctx, 2, EmptyLink)
	fmt.Println(readLink.String())
	id, _ = LinksDB.Create(ctx, readLink)
	fmt.Printf("Link with id %d is created.\n", id)
	readLink.LongLink = "http://updated.loc"
	updatedLink, _ := LinksDB.Update(ctx, readLink)
	fmt.Printf("Link with id %d is updated.\n%v\n", updatedLink.ID, updatedLink)
	LinksDB.Delete(ctx, 7)

	users := objrepo.UsersNew(UsersDB, logger)
	users.Create(ctx, &models.User{Username: "testrepo"})
	readRepoUser, _ := users.Read(ctx, 3)
	fmt.Println(readRepoUser.String())

	//users := objrepo.UsersNew(UsersDB, logger, tracer)
	//links := objrepo.LinksNew(LinksDB, logger, tracer)
}
