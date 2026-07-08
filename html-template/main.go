package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/kakkky/hotwire-go/turbo"
	"github.com/kakkky/hotwire-go/view"
)

//go:embed all:templates
var templatesFS embed.FS

func main() {
	r, err := view.New(templatesFS, "templates", view.WithFuncs(turbo.TemplateFuncMap()))
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		view:  r,
		store: newStore(),
	}

	mux := http.NewServeMux()
	s.routes(mux)

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
