package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "sync"

    "github.com/go-chi/chi/v5"
)

type Review struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

var (
    mu     sync.Mutex
    store  = map[int]Review{}
    nextID = 1
)

func main() {
    r := chi.NewRouter()

    r.Get("/health", handleHealth)
    r.Get("/reviews", handleListReviews)
    r.Get("/reviews/{id}", handleGetReview)
    r.Post("/reviews", handleCreateReview)
    r.Put("/reviews/{id}", handleUpdateReview)

    log.Println("server starting on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleListReviews(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    arr := make([]Review, 0, len(store))
    for _, v := range store {
        arr = append(arr, v)
    }
    writeJSON(w, http.StatusOK, arr)
}

func handleGetReview(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }
    mu.Lock()
    rev, ok := store[id]
    mu.Unlock()
    if !ok {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
        return
    }
    writeJSON(w, http.StatusOK, rev)
}

func handleCreateReview(w http.ResponseWriter, r *http.Request) {
    var in struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    mu.Lock()
    id := nextID
    nextID++
    rev := Review{ID: id, Title: in.Title, Content: in.Content}
    store[id] = rev
    mu.Unlock()
    writeJSON(w, http.StatusCreated, rev)
}

func handleUpdateReview(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }
    var in struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
        return
    }
    mu.Lock()
    rev, ok := store[id]
    if !ok {
        mu.Unlock()
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
        return
    }
    if in.Title != "" {
        rev.Title = in.Title
    }
    if in.Content != "" {
        rev.Content = in.Content
    }
    store[id] = rev
    mu.Unlock()
    writeJSON(w, http.StatusOK, rev)
}
