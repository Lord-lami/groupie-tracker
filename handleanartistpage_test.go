package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAnArtistPage(t *testing.T) {
	type record struct {
		id      string
		wantCode int
	}
	test := record{}
	tests := []record{}

	// Test 1 - Regular artist ID
	test.id = "9"
	test.wantCode = http.StatusOK
	tests = append(tests, test)

	// Test 2 - Excessive aritist ID
	test.id = "1000"
	test.wantCode = http.StatusNotFound
	tests = append(tests, test)

	// Test 3 - Malformed artist ID
	test.id = "abc"
	test.wantCode = http.StatusBadRequest
	tests = append(tests, test)

	// Test 4 - Empth artist ID
	test.id = ""
	test.wantCode = http.StatusBadRequest
	tests = append(tests, test)

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/artists/", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("id", test.id)

		w := httptest.NewRecorder()
		// http.DefaultServeMux.ServeHTTP(w, req)
		handleAnArtistPage(w, req)
		result := w.Result()
		// body, _ := io.ReadAll(result.Body)

		if result.StatusCode != test.wantCode {
			t.Errorf("Tested handleAnArtistPage() with url: %s", test.id)
			t.Errorf("expected %d got %d", test.wantCode, result.StatusCode)
		}
	}

}
