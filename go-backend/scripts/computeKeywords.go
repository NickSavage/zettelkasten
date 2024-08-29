package main

import (
	"database/sql"
	"fmt"
	"go-backend/models"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDatabase(dbConfig models.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v "+
		"password=%v dbname=%v sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DatabaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db, err
}

func main() {

	var s *Server
	s = &Server{}
	dbConfig := models.DatabaseConfig{}
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Password = os.Getenv("DB_PASS")
	dbConfig.DatabaseName = os.Getenv("DB_NAME")

	db, err := ConnectToDatabase(dbConfig)

	s.db = db
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err.Error())
		return
	}

	rows, _ := db.Query("SELECT id FROM users")
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Fatalf("something is wrong: %v", err.Error())
			return
		}

		log.Printf("user %v", userID)
		cards, _ := s.QueryPartialCard(userID, "")
		log.Printf("cards %v", len(cards))
	}

}
