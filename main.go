package main

import (
	"context"
	"encoding/json"
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

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ref := getBooksRef()
	var books map[string]Book

	if err := ref.Get(ctx, &books); err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrive books: %v", err), http.StatusInternalServerError)
		return
	}

	// Let's convert the map to a slice of books for a cleaner JSON output
	var bookList []Book
	for id, book := range books {
		book.ID = id
		bookList = append(bookList, book)
	}

	json.NewEncoder(w).Encode(bookList)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["id"]

	ref := getBooksRef().Child(bookID)
	var book Book
	if err := ref.Get(ctx, &book); err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrive book: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the struct has zero-values
	if book.Title == "" {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	book.ID = bookID
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// We are first creating an empty child node with an ID and getting the reference back
	ref, err := getBooksRef().Push(ctx, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusInternalServerError)
		return
	}

	// Using set, a value is stored in the database
	newBook.Added = time.Now()
	if err := ref.Set(ctx, newBook); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create book: %v", err), http.StatusInternalServerError)
		return
	}

	newBook.ID = ref.Key
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["id"]

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	ref := getBooksRef().Child(bookID)
	if err := ref.Set(ctx, updatedBook); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update book: %v", err), http.StatusInternalServerError)
		return
	}

	updatedBook.ID = bookID
	json.NewEncoder(w).Encode(updatedBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["id"]

	ref := getBooksRef().Child(bookID)
	if err := ref.Delete(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete bok: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}
