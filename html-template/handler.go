package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kakkky/hotwire-go/turbo"
	"github.com/kakkky/hotwire-go/view"
)

type server struct {
	view  *view.Renderer
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

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusFound)
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	s.view.Render(w, http.StatusOK, "index", map[string]any{
		"Todos": s.store.list(),
	})
}

func (s *server) new(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	if turbo.IsFrameRequest(r) {
		s.view.RenderPartial(w, http.StatusOK, "new_form", data)
		return
	}
	s.view.Render(w, http.StatusOK, "new", data)
}

func (s *server) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		data := map[string]any{
			"Error": "タイトルを入力してください",
			"Title": title,
		}
		if turbo.IsFrameRequest(r) {
			s.view.RenderPartial(w, http.StatusUnprocessableEntity, "new_form", data)
			return
		}
		s.view.Render(w, http.StatusUnprocessableEntity, "new", data)
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
	data := map[string]any{"Todo": t}
	if turbo.IsFrameRequest(r) {
		s.view.RenderPartial(w, http.StatusOK, "todo_edit", data)
		return
	}
	s.view.Render(w, http.StatusOK, "edit", data)
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
		data := map[string]any{
			"Todo":  &Todo{ID: id, Title: title, Done: done},
			"Error": "タイトルを入力してください",
		}
		if turbo.IsFrameRequest(r) {
			s.view.RenderPartial(w, http.StatusUnprocessableEntity, "todo_edit", data)
			return
		}
		s.view.Render(w, http.StatusUnprocessableEntity, "edit", data)
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
