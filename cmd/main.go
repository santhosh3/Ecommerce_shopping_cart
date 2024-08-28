package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/cmd/api"
)

func main() {
	//taking psqlString from ENV

	connectionString := config.Envs.PostgresString
	if len(connectionString) == 0 {
		log.Fatal("POSTGRES_SQL is not set in .env file")
	}

	//connecting to postgres DB

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec(config.Tables)

	if err != nil {
		log.Fatal(err)
	}

	//checking weather it is connected or not.
	initStorage(db)

	//running API server
	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: connected successfully")
}
