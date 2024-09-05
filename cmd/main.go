package main

import (
	"fmt"
	"log"

	"github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/cmd/api"
	"github.com/santhosh3/ECOM/database"
	"github.com/santhosh3/ECOM/models"
	"gorm.io/gorm"
)

func main() {
	//taking psqlString from ENV
	connectionString := config.Envs.PostgresString
	if len(connectionString) == 0 {
		log.Fatal("POSTGRES_SQL is not set in .env file")
	}

	//connecting to postgres DB
	db, err := database.NewPSQLStorage(connectionString)

	if err != nil {
		log.Fatal(err)
	}

	//taking RedisString from ENV
	redisString := config.Envs.RedisDB
	if len(redisString) == 0 {
		log.Fatal("POSTGRES_SQL is not set in .env file")
	}

	//connecting to redis DB
	rdb, err := database.RedisRateLimit(redisString)
     fmt.Println(rdb);

	if err != nil {
		log.Fatal(err)
	}

	// doing migrations
	models.DBMigrations(db)

	//checking DB connections
	initStorage(db)

	//running API server
	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db, rdb)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *gorm.DB) {
	sqlDB, err := db.DB() // Get the underlying sql.DB object from the GORM DB object
	if err != nil {
		log.Fatal("Failed to get database handle:", err)
	}

	// Ping the database to check if it's reachable
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("DB: connected successfully")
}
