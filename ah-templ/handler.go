package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/kakkky/hotwire-go/turbo"
)

type server struct {
	store *store
}

func (s *server) routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.root)
	mux.HandleFunc("GET /todos", s.index)
	mux.HandleFunc("GET /todos/new", s.new)
	mux.HandleFunc("POST /todos", s.create)
	mux.HandleFunc("GET /todos/{id}/edit", s.edit)
	mux.HandleFunc("POST /todos/{id}", s.update)
	mux.HandleFunc("DELETE /todos/{id}", s.delete)
}

func pathID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return 0, false
	}
	return id, true
}

func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("render: %v", err)
	}
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusFound)
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, Index(s.store.list()))
}

func (s *server) new(w http.ResponseWriter, r *http.Request) {
	if turbo.IsFrameRequest(r) {
		render(w, r, http.StatusOK, NewForm("", ""))
		return
	}
	render(w, r, http.StatusOK, New("", ""))
}

func (s *server) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		if turbo.IsFrameRequest(r) {
			render(w, r, http.StatusUnprocessableEntity, NewForm(title, "タイトルを入力してください"))
			return
		}
		render(w, r, http.StatusUnprocessableEntity, New(title, "タイトルを入力してください"))
		return
	}
	s.store.create(title)
	turbo.Redirect(w, r, "/todos")
}

func (s *server) edit(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	t, ok := s.store.get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if turbo.IsFrameRequest(r) {
		render(w, r, http.StatusOK, TodoEdit(t, ""))
		return
	}
	render(w, r, http.StatusOK, Edit(t, ""))
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	time.Sleep(2 * time.Second)
	title := r.FormValue("title")
	done := r.FormValue("done") == "on"
	if title == "" {
		t := &Todo{ID: id, Title: title, Done: done}
		if turbo.IsFrameRequest(r) {
			render(w, r, http.StatusUnprocessableEntity, TodoEdit(t, "タイトルを入力してください"))
			return
		}
		render(w, r, http.StatusUnprocessableEntity, Edit(t, "タイトルを入力してください"))
		return
	}
	if err := s.store.update(id, title, done); err != nil {
		http.NotFound(w, r)
		return
	}
	turbo.Redirect(w, r, "/todos")
}

func (s *server) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	s.store.delete(id)
	turbo.Redirect(w, r, "/todos")
}
