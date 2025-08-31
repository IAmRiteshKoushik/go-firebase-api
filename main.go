package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	firebaseConfigFile = "firebaseConfig.json"
)

type Book struct {
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   string    `json:"year"`
	Added  time.Time `json:"added"`
}

type BooksWithID struct {
	ID string `json:"id"`
	Book
}

var (
	ctx    context.Context
	app    *firebase.App
	client *firestore.Client
)

func main() {
	ctx = context.Background()
	opt := option.WithCredentialsFile(firebaseConfigFile)
	app, err := firebase.NewApp(ctx, nil, opt) // Firestore doesn't need a config
	if err != nil {
		log.Fatalf("Firebase initialization error: %v\n", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestore initialization error: %v\n", err)
	}
	defer client.Close()

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

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	iter := client.Collection("books").Documents(ctx)
	var books []BooksWithID

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to iterate document: %v", err), http.StatusInternalServerError)
			return
		}

		var book Book
		if err := doc.DataTo(&book); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse document data: %v", err), http.StatusInternalServerError)
			return
		}

		books = append(books, BooksWithID{
			ID:   doc.Ref.ID,
			Book: book,
		})
	}

	json.NewEncoder(w).Encode(books)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["id"]

	docRef := client.Collection("books").Doc(bookID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retriveb book: %v", err), http.StatusInternalServerError)
		return
	}

	var book Book
	if err := docSnap.DataTo(&book); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse document data: %v", err), http.StatusInternalServerError)
		return
	}

	bookWithID := BooksWithID{
		ID:   docSnap.Ref.ID,
		Book: book,
	}
	json.NewEncoder(w).Encode(bookWithID)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	newBook.Added = time.Now()

	docRef, _, err := client.Collection("books").Add(ctx, newBook)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create book: %v", err), http.StatusInternalServerError)
		return
	}

	bookWithID := BooksWithID{
		ID:   docRef.ID,
		Book: newBook,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookWithID)
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

	_, err := client.Collection("books").Doc(bookID).Set(ctx, updatedBook)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update book: %v", err), http.StatusInternalServerError)
		return
	}

	bookWithID := BooksWithID{ID: bookID, Book: updatedBook}
	json.NewEncoder(w).Encode(bookWithID)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["id"]

	// The ignored field here is of WriteResult type which is just a timestamp
	_, err := client.Collection("books").Doc(bookID).Delete(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete book: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}
