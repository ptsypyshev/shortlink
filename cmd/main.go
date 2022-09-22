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
	"log"

	"go.uber.org/zap"

	"github.com/ptsypyshev/shortlink/internal/app"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot initialize zap logger")
	}
	defer func() { _ = logger.Sync() }()
	//tracer, closer := InitJaeger(logger)
	//tracer, closer := InitJaeger("Shortlink", "localhost:6831", logger)
	//defer closer()

	a := app.App{}
	if err := a.Init(); err != nil {
		log.Fatalf("cannot initialize web application: %s", err)
	}
	//if closer, err := a.Init(); err != nil {
	//	log.Fatal(err)
	//} else {
	//	defer closer.Close()
	//}

	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
