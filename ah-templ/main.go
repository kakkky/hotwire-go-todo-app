package main

import (
	"log"
	"net/http"
)

func main() {
	s := &server{store: newStore()}

	mux := http.NewServeMux()
	s.routes(mux)

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
