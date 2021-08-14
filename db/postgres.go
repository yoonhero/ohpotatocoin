package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/yoonhero/ohpotatocoin/utils"
)

var DATABASE_URL string = "postgresql://yoonseonghyeon:randompassword@localhost:5432/instaclone?schema=public"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "randompassword"
	dbname   = "go_project"
)

var sqlDB *sql.DB

var once sync.Once

func dsn() string {
	dataBase := os.Getenv("DATABASE_URL")
	if dataBase != "" {
		return dataBase
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func InsertData() {
	defer CloseSqlDB()
	insertStmt := `INSERT INTO Students(Name, Roll) values('John', 1)`
	_, e := sqlDB.Exec(insertStmt)
	utils.HandleErr(e)
}

func LoadData() {
	defer CloseSqlDB()

	rows, err := sqlDB.Query(`SELECT Name, Roll FROM Students`)
	utils.HandleErr(err)

	defer rows.Close()
	for rows.Next() {
		var name string
		var roll int

		err = rows.Scan(&name, &roll)
		utils.HandleErr(err)

		fmt.Println(name, roll)
	}

}

func createTable() {
	// blocks table and data table
	stmt, err := sqlDB.Prepare("CREATE TABLE Students (Name varchar(111) NOT NULL, Roll int)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)

	fmt.Println("Table created successfully..")
}

func CloseSqlDB() {
	sqlDB.Close()
}

func InitPostgresDB() {
	if sqlDB == nil {
		db, err := sql.Open("postgres", dsn())
		utils.HandleErr(err)

		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		err = db.PingContext(ctx)
		utils.HandleErr(err)
		log.Printf("Connected to DB %s successfully\n", dbname)

		sqlDB = db

		createTable()

	}
}

// save block data
func saveBlockInSQL(hash string, data []byte) {
	// update database
	defer CloseSqlDB()
	insertStmt := fmt.Sprintf("INSERT INTO Blocks(Hash, Data) values(%s, %x)", hash, data)
	_, err := sqlDB.Exec(insertStmt)
	utils.HandleErr(err)
}

// empty chain table
func emptyChainTable() {
	defer CloseSqlDB()
}

// save chain
func saveChainInSQL(data []byte) {
	defer CloseSqlDB()
	insertCommand := fmt.Sprintf("INSERT INTO Checkpoint(Data) values(%v)", data)
	_, err := sqlDB.Exec(insertCommand)
	utils.HandleErr(err)
}
