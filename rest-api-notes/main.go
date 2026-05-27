package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Store struct {
	mu     sync.RWMutex
	notes  map[int]Note
	nextID int
}

func NewStore() *Store {
	return &Store{
		notes:  make(map[int]Note),
		nextID: 1,
	}
}

func (s *Store) Add(title, body string) Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	note := Note{ID: s.nextID, Title: title, Body: body}
	s.notes[s.nextID] = note
	s.nextID++
	return note
}

func (s *Store) GetAll() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes := make([]Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes
}

func (s *Store) Search(query string, sort string, limit int) []Note {	
	s.mu.RLock()
	defer s.mu.RUnlock()
	query = strings.ToLower(query)
	results := []Note{}
	for _, note := range s.notes {
		if strings.Contains(strings.ToLower(note.Title), query) ||
			strings.Contains(strings.ToLower(note.Body), query) {
			results = append(results, note)
		}
	}

	switch sort {
		case "asc":
			slices.SortFunc(results, func(a, b Note) int {
				return a.ID - b.ID
			})
		case "desc":
			slices.SortFunc(results, func(a, b Note) int {
				return b.ID - a.ID
			})
	}

	if len(results) > limit {
		results = results[:limit]
	}

	if len(results) > limit {
		results = results[:limit]
	}
	return results
}

func (s *Store) Get(id int) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, ok := s.notes[id]
	return note, ok
}

func (s *Store) Update(id int, title, body string) (Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	note, ok := s.notes[id]
	if !ok {
		return Note{}, false
	}
	note.Title = title
	note.Body = body
	s.notes[id] = note
	return note, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.notes[id]
	if ok {
		delete(s.notes, id)
	}
	return ok
}

// --- middleware ---

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s → %d (%s)",
			r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

func jsonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// --- helpers ---

func sortNotes(notes []Note, ) {
	
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func idFromPath(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	return strconv.Atoi(parts[len(parts)-1])
}

// --- handlers ---

func (s *Store) listNotes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if query != "" {
		writeJSON(w, http.StatusOK, s.Search(query, sort, limit))
		return
	}
	writeJSON(w, http.StatusOK, s.GetAll())
}

func (s *Store) createNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	note := s.Add(input.Title, input.Body)
	writeJSON(w, http.StatusCreated, note)
}

func (s *Store) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	note, ok := s.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (s *Store) updateNote(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	note, ok := s.Update(id, input.Title, input.Body)
	if !ok {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (s *Store) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if !s.Delete(id) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func main() {
	store := NewStore()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", store.listNotes)
	mux.HandleFunc("POST /notes", store.createNote)
	mux.HandleFunc("GET /notes/{id}", store.getNote)
	mux.HandleFunc("PUT /notes/{id}", store.updateNote)
	mux.HandleFunc("DELETE /notes/{id}", store.deleteNote)

	handler := chain(mux, logger, jsonHeaders)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}