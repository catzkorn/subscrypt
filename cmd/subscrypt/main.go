package main

import (
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"log"
	"net/http"
	"os"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/server"
)

func main() {

	database, err := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	if err != nil {
		log.Fatalf("failed to create database connection: %v", err)
	}

	api := &plaid.PlaidAPI{}

	server := server.NewServer(database, api)
	err = http.ListenAndServe(":5000", server)
	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
