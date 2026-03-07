package reddit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRedditClientValidation(t *testing.T) {
	tests := []struct {
		name         string
		clientId     string
		clientSecret string
		username     string
		password     string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "missing client id",
			clientId:     "",
			clientSecret: "test-client-secret",
			username:     "test-user",
			password:     "test-password",
			wantErr:      true,
			errContains:  "required",
		},
		{
			name:         "missing client secret",
			clientId:     "test-client-id",
			clientSecret: "",
			username:     "test-user",
			password:     "test-password",
			wantErr:      true,
			errContains:  "required",
		},
		{
			name:         "missing username",
			clientId:     "test-client-id",
			clientSecret: "test-client-secret",
			username:     "",
			password:     "test-password",
			wantErr:      true,
			errContains:  "required",
		},
		{
			name:         "missing password",
			clientId:     "test-client-id",
			clientSecret: "test-client-secret",
			username:     "test-user",
			password:     "",
			wantErr:      true,
			errContains:  "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedditClient(tt.clientId, tt.clientSecret, tt.username, tt.password, "", 0, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedditClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("NewRedditClient() error = %v, should contain %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestRetryHttpRequest(t *testing.T) {
	tests := []struct {
		name        string
		statusCodes []int
		attempts    int
		wantSuccess bool
	}{
		{
			name:        "success on first attempt",
			statusCodes: []int{http.StatusOK},
			attempts:    3,
			wantSuccess: true,
		},
		{
			name:        "eventual success after failures",
			statusCodes: []int{http.StatusInternalServerError, http.StatusInternalServerError, http.StatusOK},
			attempts:    3,
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				idx := callCount
				if idx >= len(tt.statusCodes) {
					idx = len(tt.statusCodes) - 1
				}
				w.WriteHeader(tt.statusCodes[idx])
				callCount++
			}))
			defer server.Close()

			req, _ := http.NewRequest("GET", server.URL, nil)
			client := &http.Client{}

			// Use a short delay for testing
			resp, err := retryHttpRequest(client, req, tt.attempts, 1*time.Millisecond)

			if tt.wantSuccess {
				if err != nil {
					t.Errorf("retryHttpRequest() unexpected error = %v", err)
				}
				if resp == nil || resp.StatusCode/100 != 2 {
					t.Error("retryHttpRequest() expected successful response")
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
