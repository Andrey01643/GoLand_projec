package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubstringHandler(t *testing.T) {
	testCases := []struct {
		input       string
		expected    string
		expectedErr string
	}{
		{"abcabcbb", "abc", ""},
		{"bbbbb", "b", ""},
		{"pwwkew", "wke", ""},
		{"pww kew", "wke", "Empty string"},
		{"", "", "Empty string"},
		{" ", "", "Empty string"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", "/?str="+tc.input, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SubstringHandler)
		handler.ServeHTTP(rr, req)

		if tc.expectedErr != "" {
			if rr.Code != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedErr) {
				t.Errorf("Expected error: %s, but got: %s", tc.expectedErr, rr.Body.String())
			}
		} else {
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != tc.expected {
				t.Errorf("Input: %s, Expected: %s, Got: %s", tc.input, tc.expected, strings.TrimSpace(rr.Body.String()))
			}
		}
	}
}
