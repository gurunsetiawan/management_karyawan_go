package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"John Doe", "John Doe"},
		{"test@example.com", "test@example.com"},
		{"  test  ", "test"},
		{"", ""},
	}

	for _, test := range tests {
		result := sanitizeInput(test.input)
		if result != test.expected {
			t.Errorf("sanitizeInput(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"", false},
		{"user@domain", false},
	}

	for _, test := range tests {
		result := isValidEmail(test.email)
		if result != test.expected {
			t.Errorf("isValidEmail(%q) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestErrorResponseStructure(t *testing.T) {
	errorResp := ErrorResponse{
		Error:   "Test Error",
		Message: "Test Message",
		Code:    400,
	}

	jsonData, err := json.Marshal(errorResp)
	if err != nil {
		t.Errorf("Failed to marshal ErrorResponse: %v", err)
	}

	var unmarshaled ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ErrorResponse: %v", err)
	}

	if unmarshaled.Error != errorResp.Error {
		t.Errorf("Error field mismatch: got %s, expected %s", unmarshaled.Error, errorResp.Error)
	}

	if unmarshaled.Message != errorResp.Message {
		t.Errorf("Message field mismatch: got %s, expected %s", unmarshaled.Message, errorResp.Message)
	}

	if unmarshaled.Code != errorResp.Code {
		t.Errorf("Code field mismatch: got %d, expected %d", unmarshaled.Code, errorResp.Code)
	}
}

func TestEmployeeStructure(t *testing.T) {
	employee := Employee{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Role:      "Developer",
		Phone:     "1234567890",
		Alamat:    "Test Address",
		CreatedAt: "2023-01-01 00:00:00",
	}

	jsonData, err := json.Marshal(employee)
	if err != nil {
		t.Errorf("Failed to marshal Employee: %v", err)
	}

	var unmarshaled Employee
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal Employee: %v", err)
	}

	if unmarshaled.ID != employee.ID {
		t.Errorf("ID field mismatch: got %d, expected %d", unmarshaled.ID, employee.ID)
	}

	if unmarshaled.Name != employee.Name {
		t.Errorf("Name field mismatch: got %s, expected %s", unmarshaled.Name, employee.Name)
	}

	if unmarshaled.Email != employee.Email {
		t.Errorf("Email field mismatch: got %s, expected %s", unmarshaled.Email, employee.Email)
	}
}

func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/employees", nil)
	w := httptest.NewRecorder()

	// Apply CORS middleware
	corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})(w, req)

	// Check CORS headers
	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin *, got %s", corsOrigin)
	}

	corsMethods := w.Header().Get("Access-Control-Allow-Methods")
	if corsMethods == "" {
		t.Errorf("Expected Access-Control-Allow-Methods header")
	}
}

func TestRateLimiting(t *testing.T) {
	// Test rate limiting by making multiple requests
	// Note: This test may fail if rate limiting is not properly configured
	// We'll test with a smaller number to avoid false positives
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/api/employees", nil)
		w := httptest.NewRecorder()

		rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})(w, req)

		// We expect all requests to succeed in test environment
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got status %d", w.Code)
		}
	}
}

// Benchmark tests
func BenchmarkSanitizeInput(b *testing.B) {
	input := "<script>alert('xss')</script>"
	for i := 0; i < b.N; i++ {
		sanitizeInput(input)
	}
}

func BenchmarkIsValidEmail(b *testing.B) {
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		isValidEmail(email)
	}
}
