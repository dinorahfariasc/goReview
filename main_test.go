package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func resetStore() {
    mu.Lock()
    store = map[int]Review{}
    nextID = 1
    mu.Unlock()
}

func TestHealth(t *testing.T) {
    router := newRouter()
    req := httptest.NewRequest("GET", "/health", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200 got %d", rr.Code)
    }
    var body map[string]string
    if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
        t.Fatalf("decode failed: %v", err)
    }
    if body["status"] != "ok" {
        t.Fatalf("unexpected body: %v", body)
    }
}

func TestCreateGetUpdateList(t *testing.T) {
    resetStore()
    router := newRouter()

    // Create
    createBody := bytes.NewBufferString(`{"title":"T","content":"C"}`)
    req := httptest.NewRequest("POST", "/reviews", createBody)
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    if rr.Code != http.StatusCreated {
        t.Fatalf("create expected 201 got %d", rr.Code)
    }
    var created Review
    if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
        t.Fatalf("decode created: %v", err)
    }
    if created.ID != 1 || created.Title != "T" {
        t.Fatalf("created mismatch: %v", created)
    }

    // Get
    req = httptest.NewRequest("GET", "/reviews/1", nil)
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("get expected 200 got %d", rr.Code)
    }
    var got Review
    if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
        t.Fatalf("decode get: %v", err)
    }
    if got.ID != 1 || got.Content != "C" {
        t.Fatalf("get mismatch: %v", got)
    }

    // Update
    updateBody := bytes.NewBufferString(`{"content":"Updated"}`)
    req = httptest.NewRequest("PUT", "/reviews/1", updateBody)
    req.Header.Set("Content-Type", "application/json")
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("update expected 200 got %d", rr.Code)
    }
    var updated Review
    if err := json.NewDecoder(rr.Body).Decode(&updated); err != nil {
        t.Fatalf("decode update: %v", err)
    }
    if updated.Content != "Updated" {
        t.Fatalf("update failed: %v", updated)
    }

    // List
    req = httptest.NewRequest("GET", "/reviews", nil)
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("list expected 200 got %d", rr.Code)
    }
    var list []Review
    if err := json.NewDecoder(rr.Body).Decode(&list); err != nil {
        t.Fatalf("decode list: %v", err)
    }
    if len(list) != 1 {
        t.Fatalf("expected list len 1 got %d", len(list))
    }
}
