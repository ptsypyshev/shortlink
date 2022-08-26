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

	id := 3
	UsersDB := pgdb.DBNew[*models.User](pool)
	users := objrepo.UsersNew(UsersDB, logger)
	createdUser, _ := users.Create(ctx, &models.User{Username: "testrepo"})
	fmt.Printf("New user created: %v", createdUser)
	readRepoUser, _ := users.Read(ctx, id)
	fmt.Printf("User with id %d is read: %v", id, readRepoUser)
	updatedUser, _ := users.Update(ctx, id, &models.User{FirstName: "Alex"})
	fmt.Printf("User with id %d is updated: %v", id, updatedUser)
	id++
	deletedUser, _ := users.Delete(ctx, id)
	fmt.Printf("User with id %d is deleted: %v", id, deletedUser)
	foundUsers, _ := users.Search(ctx, "username", "admin")
	fmt.Printf("Users with username %v are found: %v", id, foundUsers)

	id = 2
	LinksDB := pgdb.DBNew[*models.Link](pool)
	links := objrepo.LinksNew(LinksDB, logger)
	createdLink, _ := links.Create(ctx, &models.Link{LongLink: "https://yahoo.com", OwnerID: 3})
	fmt.Printf("New link created: %v", createdLink)
	readRepoLink, _ := links.Read(ctx, id)
	fmt.Printf("Link with id %d is read: %v", id, readRepoLink)
	updatedLink, _ := links.Update(ctx, id, &models.Link{LongLink: "http://localhost:8080", OwnerID: 2})
	fmt.Printf("Link with id %d is updated: %v", id, updatedLink)
	id++
	deletedLink, _ := links.Delete(ctx, id)
	fmt.Printf("Link with id %d is deleted: %v", id, deletedLink)
	foundLinks, _ := links.Search(ctx, "owner_id", id)
	fmt.Printf("Links with owner id %v are found: %v", id, foundLinks)

	id = 4
	ShortLinksDB := pgdb.DBNew[*models.ShortLink](pool)
	shortlinks := objrepo.ShortLinksNew(ShortLinksDB, logger)
	createdShortLink, _ := shortlinks.Create(ctx, id)
	fmt.Printf("New shortlink created: %v", createdShortLink)
	readRepoShortLink, _ := shortlinks.Read(ctx, id)
	fmt.Printf("ShortLink with id %d is read: %v", id, readRepoShortLink)
	updatedShortLink, _ := shortlinks.Update(ctx, id, &models.ShortLink{Token: "6fanhhaat0", LongLinkID: 1})
	fmt.Printf("ShortLink with id %d is updated: %v", id, updatedShortLink)
	id++
	deletedShortLink, _ := shortlinks.Delete(ctx, id)
	fmt.Printf("ShortLink with id %d is deleted: %v", id, deletedShortLink)
	foundShortLinks, _ := shortlinks.Search(ctx, "token", "tr17bwgh9")
	fmt.Printf("Shortlink with token tr17bwgh9 are found: %v", foundShortLinks)

	//for i := 0; i < 10; i++ {
	//	token, err := objrepo.GenerateShortLinkToken(i)
	//	if err != nil {
	//		fmt.Printf("cannot generate token for id %d: %s", i, err)
	//	}
	//	fmt.Printf("Hash for id %d is %s\n", i, token)
	//}
}
