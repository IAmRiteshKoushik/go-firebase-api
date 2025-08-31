# Bookstore CRUD App (Golang + Firebase)

This is a simple CRUD (Create, Read, Update, Delete) application for 
managing a bookstore, built with **Golang** as the backend and **Firebase** 
as the data store. The goal is to test and explore integrating Firebase and Go.

## Features

- Add new books (Create)
- View a list or details of books (Read)
- Update book details (Update)
- Remove books from the store (Delete)

## Prerequisites

- Go 1.23+ installed
- Firebase project with Firestore database enabled
- Service Account JSON from Firebase Console

## Setup

1. **Clone the Repository**

    ```
    git clone https://github.com/yourusername/go-firebase-bookstore.git
    cd go-firebase-bookstore
    ```
2. **Install Go Dependencies**
   
    ```
    go mod tidy
    ```
3. **Add Firebase Credentials**

    - Download your Firebase project’s service account JSON.
    - Save it locally (e.g., `serviceAccountKey.json`).

4. **Set GOOGLE_APPLICATION_CREDENTIALS**

    ```
    export GOOGLE_APPLICATION_CREDENTIALS="path/to/serviceAccountKey.json"
    ```

5. **Run the Application**

    ```
    go run main.go
    ```

## Project Structure

- `main.go` — entry point; HTTP handlers for CRUD operations

## Example API Endpoints

- `POST   /api/books` — Create a new book
- `GET    /api/books` — List all books
- `GET    /api/books/{id}` — Get details for a book
- `PUT    /api/books/{id}` — Update book info
- `DELETE /api/books/{id}` — Delete a book

---

References:
- [Go Quickstart for Firestore](https://firebase.google.com/docs/fire
