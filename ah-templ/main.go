package main

import (
	"log"
	"net/http"

	"github.com/kakkky/hotwire-go/turbo"
)

func main() {
	sb := turbo.NewStreamBroker()
	s := &server{store: newStore(), broker: sb}

	mux := http.NewServeMux()
	mux.Handle(turbo.StreamsSSEPath, turbo.StreamSSEHandler(sb))

	pageMux := http.NewServeMux()
	s.routes(pageMux)
	mux.Handle("/", turbo.StreamsMiddleware(pageMux))

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
