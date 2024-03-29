package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"webapp/pkg/repository"
)

const port = 8080

type application struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app application

	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application, e.g. company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres Connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "eraser-secret", "signing secret")
	flag.Parse()

	// connect to DB
	// conn, err := app.connectToDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	// app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Printf("Starting api on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
