package main

import (
	"log"

	"github.com/reyhardy/exs-prescription-app/internal/db"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := db.New("sqlite", "app.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	if err := run(db); err != nil {
		log.Fatal(err)
	}
}
