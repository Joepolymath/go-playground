package web

import (
	"log"
	"net/http"
)

const PORT = ":4001"

func Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets", displaySnippet)
	mux.HandleFunc("/snippets/create", createSnippet)

	log.Printf("Starting server on port %s", PORT)
	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatal(err)
	}
}
