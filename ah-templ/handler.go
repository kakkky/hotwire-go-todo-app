package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusFound)
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := Index(s.store.list()).Render(r.Context(), w); err != nil {
		log.Printf("index: %v", err)
	}
}

func (s *server) new(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if turbo.IsFrameRequest(r) {
		if err := NewForm("", "").Render(r.Context(), w); err != nil {
			log.Printf("new: %v", err)
		}
		return
	}
	if err := New("", "").Render(r.Context(), w); err != nil {
		log.Printf("new: %v", err)
	}
}

func (s *server) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		if turbo.IsStreamRequest(r) {
			turbo.StreamHeader(w)
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := CreateFailStreams(title, "タイトルを入力してください").Render(r.Context(), w); err != nil {
				log.Printf("create: %v", err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if turbo.IsFrameRequest(r) {
			if err := NewForm(title, "タイトルを入力してください").Render(r.Context(), w); err != nil {
				log.Printf("create: %v", err)
			}
			return
		}
		if err := New(title, "タイトルを入力してください").Render(r.Context(), w); err != nil {
			log.Printf("create: %v", err)
		}
		return
	}
	t := s.store.create(title)
	if turbo.IsStreamRequest(r) {
		turbo.StreamHeader(w)
		w.WriteHeader(http.StatusOK)
		if err := CreateSuccessStreams(t).Render(r.Context(), w); err != nil {
			log.Printf("create: %v", err)
		}
		return
	}
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if turbo.IsFrameRequest(r) {
		if err := TodoEdit(t, "").Render(r.Context(), w); err != nil {
			log.Printf("edit: %v", err)
		}
		return
	}
	if err := Edit(t, "").Render(r.Context(), w); err != nil {
		log.Printf("edit: %v", err)
	}
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
		if turbo.IsStreamRequest(r) {
			turbo.StreamHeader(w)
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := UpdateFailStreams(t, "タイトルを入力してください").Render(r.Context(), w); err != nil {
				log.Printf("update: %v", err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if turbo.IsFrameRequest(r) {
			if err := TodoEdit(t, "タイトルを入力してください").Render(r.Context(), w); err != nil {
				log.Printf("update: %v", err)
			}
			return
		}
		if err := Edit(t, "タイトルを入力してください").Render(r.Context(), w); err != nil {
			log.Printf("update: %v", err)
		}
		return
	}
	if err := s.store.update(id, title, done); err != nil {
		http.NotFound(w, r)
		return
	}
	if turbo.IsStreamRequest(r) {
		t, _ := s.store.get(id)
		turbo.StreamHeader(w)
		w.WriteHeader(http.StatusOK)
		if err := UpdateSuccessStreams(t).Render(r.Context(), w); err != nil {
			log.Printf("update: %v", err)
		}
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
	if turbo.IsStreamRequest(r) {
		turbo.StreamHeader(w)
		w.WriteHeader(http.StatusOK)
		if err := DeleteStreams(id).Render(r.Context(), w); err != nil {
			log.Printf("delete: %v", err)
		}
		return
	}
	turbo.Redirect(w, r, "/todos")
}
