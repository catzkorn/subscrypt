package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/subscription"
)

func main() {

	database, _ := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))

	server := subscription.NewSubscriptionServer(database)
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
