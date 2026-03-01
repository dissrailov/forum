package main

import (
	"forum/internal/ai"
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	_ = godotenv.Load()

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

	aiClient := ai.NewClient()
	aiService := ai.NewService(aiClient, db, errorlog)

	if aiClient.IsConfigured() {
		infolog.Println("AI (Groq) is configured and enabled")
	} else {
		infolog.Println("AI (Groq) is not configured — set GROQ_API_KEY to enable")
	}

	service := service.NewService(db)
	app := app.New(infolog, errorlog, templateCache)
	handlers := handlers.New(service, app, aiService)

	infolog.Println("Server is running on :http://localhost:8081")
	srv := &http.Server{
		Addr:     ":8081",
		ErrorLog: app.ErrorLog,
		Handler:  handlers.Routes(),
	}
	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}
