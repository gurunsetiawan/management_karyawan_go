package handler

import (
	"net/http"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		xRealIP        string
		expected       string
	}{
		{
			name:       "Only RemoteAddr with port",
			remoteAddr: "192.168.1.10:45231",
			expected:   "192.168.1.10",
		},
		{
			name:       "Only RemoteAddr without port (should handle gracefully)",
			remoteAddr: "10.0.0.1",
			expected:   "10.0.0.1",
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[2001:db8::1]:8080",
			expected:   "2001:db8::1",
		},
		{
			name:          "With X-Forwarded-For single IP",
			remoteAddr:    "127.0.0.1:1234",
			xForwardedFor: "203.0.113.195",
			expected:      "203.0.113.195",
		},
		{
			name:          "With X-Forwarded-For multiple IPs",
			remoteAddr:    "127.0.0.1:1234",
			xForwardedFor: "203.0.113.195, 70.41.3.18, 150.172.238.178",
			expected:      "203.0.113.195", // Should take the first one (original client)
		},
		{
			name:       "With X-Real-IP",
			remoteAddr: "127.0.0.1:1234",
			xRealIP:    "198.51.100.1",
			expected:   "198.51.100.1",
		},
		{
			name:          "X-Forwarded-For takes precedence over X-Real-IP",
			remoteAddr:    "127.0.0.1:1234",
			xForwardedFor: "203.0.113.195",
			xRealIP:       "198.51.100.1",
			expected:      "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			ip := getClientIP(req)
			if ip != tt.expected {
				t.Errorf("getClientIP() = %v, want %v", ip, tt.expected)
			}
		})
	}
}

func TestParseAllowedOrigins(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		expected []string
	}{
		{
			name:     "Empty env var returns defaults",
			envVar:   "",
			expected: []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "http://localhost:8083", "http://127.0.0.1:8083"},
		},
		{
			name:     "Single origin",
			envVar:   "https://app.com",
			expected: []string{"https://app.com"},
		},
		{
			name:     "Multiple origins with spaces",
			envVar:   "https://app.com, http://test.com ,  http://localhost:3000",
			expected: []string{"https://app.com", "http://test.com", "http://localhost:3000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAllowedOrigins(tt.envVar)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Expected %s at index %d, got %s", tt.expected[i], i, v)
				}
			}
		})
	}
}
