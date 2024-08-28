package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/cmd/api"
)



func main()  {
	connectionString := config.Envs.PostgresString
	if len(connectionString) == 0 {
		log.Fatal("POSTGRES_SQL is not set in .env file")
	}
	db, err := sql.Open("postgres", connectionString)
     if err != nil {
        log.Fatal(err)
     }

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s",config.Envs.Port),db);
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB)  {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: connected successfully")
}