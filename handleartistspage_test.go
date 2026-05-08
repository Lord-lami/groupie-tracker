package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleArtistsPage(t *testing.T) {
	type record struct {
		url      string
		wantCode int
	}
	test := record{}
	tests := []record{}

	// Test 1 - Normal Base Path
	test.url = "/"
	test.wantCode = http.StatusSeeOther
	tests = append(tests, test)

	// Test 2 - Artists Path
	test.url = "/artists"
	test.wantCode = http.StatusSeeOther
	tests = append(tests, test)

	// Test 3 - Normal Artist Page Path
	test.url = "/artist?page=1"
	test.wantCode = http.StatusOK
	tests = append(tests, test)

	// Test 4 - Too long page number for Artist Path
	test.url = "/artist?page=1000"
	test.wantCode = http.StatusNotFound
	tests = append(tests, test)

	// Test 5 - Malformed page number for Artist Page Path
	test.url = "/artist?page=abc"
	test.wantCode = http.StatusBadRequest
	tests = append(tests, test)

	// Test 6 - Empty page number for Artist Page Path
	test.url = "/artist?page="
	test.wantCode = http.StatusSeeOther
	tests = append(tests, test)

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.url, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		handleArtistsPage(w, req)
		result := w.Result()
		// body, _ := io.ReadAll(result.Body)

		if result.StatusCode != test.wantCode {
			t.Errorf("Tested handleArtistsPage() with url: %s", test.url)
			t.Errorf("expected %d got %d", test.wantCode, result.StatusCode)
		}
	}

}
