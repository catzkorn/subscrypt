package main

import (
	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"log"
	"net/http"
)

func main() {

	database, _ := database.NewDatabaseConnection(database.DatabaseConnTestString)

	server := subscription.NewSubscriptionServer(database)
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
