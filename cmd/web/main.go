package main

import (
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := repo.NewRepo("./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}
	service := service.NewService(db)
	app := app.New(infolog, errorlog, templateCache)
	handlers := handlers.New(service, app)

	infolog.Println("Server is running on :http://localhost:8080")
	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: app.ErrorLog,
		Handler:  handlers.Routes(),
	}
	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}
