package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

const (
	firebaseConfigFile = "firebaseConfig.json"
	firebaseDBURL      = ""
)

type Book struct {
	ID     string    `json:"id,omitempty"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   string    `json:"year"`
	Added  time.Time `json:"added,omitempty"`
}

var (
	ctx    context.Context
	app    *firebase.App
	client *db.Client
)

func main() {
	ctx = context.Background()
	opt := option.WithCredentialsFile(firebaseConfigFile)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: firebaseDBURL,
	}, opt)
	if err != nil {
		log.Fatalf("Firebase initialization error: %v\n", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("Firestore initialization error: %v\n", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBookById).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func getBooksRef() *db.Ref {
	return client.NewRef("books")
}

func getBooks(w http.ResponseWriter, r *http.Request) {}

func getBookById(w http.ResponseWriter, r *http.Request) {}

func createBook(w http.ResponseWriter, r *http.Request) {}

func updateBook(w http.ResponseWriter, r *http.Request) {}

func deleteBook(w http.ResponseWriter, r *http.Request) {}
