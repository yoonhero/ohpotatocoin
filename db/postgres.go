package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/yoonhero/ohpotatocoin/utils"
)

var (
	host     = "ec2-18-214-208-89.compute-1.amazonaws.com"
	port     = 5434
	user     = "postgres"
	password = "cf035f3144d8e022e61f29cfeb0ea198344f886fd18763143aab0df0d2cd099d"
	dbname   = "d7nafmdveejv65"
)

var sqlDB *sql.DB

func dsn() string {
	// postgres://hwbdgkrceyodzs:cf035f3144d8e022e61f29cfeb0ea198344f886fd18763143aab0df0d2cd099d@ec2-18-214-208-89.compute-1.amazonaws.com:5432/d7nafmdveejv65

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func createCheckpointTable() {
	stmt, err := sqlDB.Prepare("CREATE TABLE IF NOT EXISTS Checkpoint (Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
}

func createBlocksTable() {
	stmt, err := sqlDB.Prepare("CREATE TABLE IF NOT EXISTS Blocks (Hash varchar(111) NOT NULL, Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
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

		createBlocksTable()
		createCheckpointTable()
	}
}

// save block data
func saveBlockInSQL(hash string, data []byte) {
	// update database
	_, err := sqlDB.Exec("INSERT INTO Blocks(Hash, Data) values($1, $2)", hash, data)
	utils.HandleErr(err)
}

// empty chain table
func emptyChainTable() {
	stmt, err := sqlDB.Prepare("DROP TABLE Checkpoint")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createCheckpointTable()
}

// save chain
func saveChainInSQL(data []byte) {
	emptyChainTable()

	_, err := sqlDB.Exec("INSERT INTO Checkpoint(Data) values($1)", data)
	utils.HandleErr(err)
}

func loadChainInSQL() []byte {
	var data []byte

	rows, err := sqlDB.Query("SELECT Data FROM Checkpoint")
	utils.HandleErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		utils.HandleErr(err)
	}

	return data
}

func findBlockInSQL(hash string) []byte {
	var data []byte

	err := sqlDB.QueryRow("SELECT Data FROM Blocks WHERE Hash = $1", hash).Scan(&data)
	utils.HandleErr(err)

	return data
}

func emptyBlocksInSQL() {
	stmt, err := sqlDB.Prepare("DROP TABLE Blocks")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createBlocksTable()
}
