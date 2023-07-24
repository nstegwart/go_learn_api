package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/go_learn_api/models"
	_ "github.com/go-sql-driver/mysql"
)

// Fungsi untuk membuat koneksi ke database MySQL
func createConnection() (*sql.DB, error) {
	// Konfigurasi koneksi ke MySQL
	username := "root"
	password := ""
	host := "localhost"
	port := "3306"
	database := "go_learn"

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Handler untuk mendapatkan semua buku
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	db, err := createConnection()
	if err != nil {
		http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, "Could not get books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			http.Error(w, "Error scanning books", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Handler untuk mendapatkan buku berdasarkan ID
func GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	db, err := createConnection()
	if err != nil {
		http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Dapatkan ID buku dari query parameter
	id := r.URL.Query().Get("id")

	var book models.Book
	err = db.QueryRow("SELECT * FROM books WHERE id=?", id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Handler untuk menambahkan buku baru
func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	db, err := createConnection()
	if err != nil {
		http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newBook models.Book
	err = json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", newBook.Title, newBook.Author)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler untuk menghapus buku berdasarkan ID
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	db, err := createConnection()
	if err != nil {
		http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Dapatkan ID buku dari query parameter
	id := r.URL.Query().Get("id")

	_, err = db.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
