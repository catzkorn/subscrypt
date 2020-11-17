package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Catzkorn/subscrypt/internal/database"
	"github.com/Catzkorn/subscrypt/internal/server"
	"github.com/sendgrid/sendgrid-go"
)

func main() {

	database, err := database.NewDatabaseConnection(os.Getenv("DATABASE_CONN_STRING"))
	if err != nil {
		log.Fatalf("failed to create database connection: %v", err)
	}

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	server := server.NewServer(database, "./web/index.html", client)
	err = http.ListenAndServe(":5000", server)
	if err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}

}
