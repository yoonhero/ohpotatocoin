package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DATABASE_URL string = "postgresql://yoonseonghyeon:randompassword@localhost:5432/instaclone?schema=public"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "randompassword"
	dbname   = "go_project"
)

func Start() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
