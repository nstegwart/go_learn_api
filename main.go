package main

import (
	"github.com/nstegwart/go_learn_api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Inisialisasi router
	router := mux.NewRouter()

	// Definisikan rute untuk endpoint API
	router.HandleFunc("/api/books", handlers.GetBooksHandler).Methods("GET")
	router.HandleFunc("/api/book", handlers.GetBookByIDHandler).Methods("GET")
	router.HandleFunc("/api/books", handlers.CreateBookHandler).Methods("POST")
	router.HandleFunc("/api/book", handlers.DeleteBookHandler).Methods("DELETE")

	// Mulai server pada port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}
