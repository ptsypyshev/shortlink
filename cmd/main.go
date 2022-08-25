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

	"github.com/ptsypyshev/shortlink/internal/db/pgdb"
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

	pool, err := pgdb.InitDB(ctx, logger)
	if err != nil {
		log.Fatalf("cannot init DB: %s", err)
	}

	if err := pgdb.InitSchema(ctx, pool); err != nil {
		panic("cannot init schema")
	}
	if err := pgdb.AddDemoData(ctx, pool); err != nil {
		panic("cannot init schema")
	}

	UsersDB := pgdb.DBNew[*models.User](pool)
	//EmptyUser := &models.User{}
	//readUser, _ := UsersDB.Read(ctx, 2, EmptyUser)
	//fmt.Println(readUser.String())
	//id, _ := UsersDB.Create(ctx, readUser)
	//fmt.Printf("User with id %d is created.\n", id)
	//readUser.Username = "Updated"
	//updatedUser, _ := UsersDB.Update(ctx, readUser)
	//fmt.Printf("User with id %d is updated.\n%v\n", updatedUser.ID, updatedUser)
	//UsersDB.Delete(ctx, 5)

	users := objrepo.UsersNew(UsersDB, logger)
	createdUser, _ := users.Create(ctx, &models.User{Username: "testrepo"})
	fmt.Printf("New user created: %v", createdUser)
	id := 3
	readRepoUser, _ := users.Read(ctx, id)
	fmt.Printf("User with id %d is read: %v", id, readRepoUser)
	updatedUser, _ := users.Update(ctx, id, &models.User{FirstName: "Alex"})
	fmt.Printf("User with id %d is updated: %v", id, updatedUser)
	id++
	deleteddUser, _ := users.Delete(ctx, id)
	fmt.Printf("User with id %d is deleted: %v", id, deleteddUser)

	LinksDB := pgdb.DBNew[*models.Link](pool)
	//EmptyLink := &models.Link{}
	//readLink, _ := LinksDB.Read(ctx, 5, EmptyLink)
	//fmt.Println(readLink.String())
	//id, _ = LinksDB.Create(ctx, readLink)
	//fmt.Printf("Link with id %d is created.\n", id)
	//readLink.LongLink = "http://updated.loc"
	//updatedLink, _ := LinksDB.Update(ctx, readLink)
	//fmt.Printf("Link with id %d is updated.\n%v\n", updatedLink.ID, updatedLink)
	//LinksDB.Delete(ctx, 7)

	links := objrepo.LinksNew(LinksDB, logger)
	createdLink, _ := links.Create(ctx, &models.Link{LongLink: "http://example.com", ShortLink: "33sfg20sn", OwnerID: 5})
	fmt.Printf("New link created: %v", createdLink)
	id = 5
	readRepoLink, _ := links.Read(ctx, id)
	fmt.Printf("Link with id %d is read: %v", id, readRepoLink)
	updatedLink, _ := links.Update(ctx, id, &models.Link{LongLink: "http://localhost:8080", OwnerID: 2})
	fmt.Printf("Link with id %d is updated: %v", id, updatedLink)
	id++
	deletedLink, _ := links.Delete(ctx, id)
	fmt.Printf("Link with id %d is deleted: %v", id, deletedLink)

	//users := objrepo.UsersNew(UsersDB, logger, tracer)
	//links := objrepo.LinksNew(LinksDB, logger, tracer)
}
