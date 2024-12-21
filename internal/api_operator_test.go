package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostForAuth(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse map[string]interface{}
		serverStatus   int
		expectedAuth   Auth
		expectedErr    bool
	}{
		{
			name: "Successful Authentication",
			serverResponse: map[string]interface{}{
				"access_token": "test-token",
				"token_type":   "bearer",
				"expires_in":   3600,
			},
			serverStatus: http.StatusOK,
			expectedAuth: Auth{
				AccessToken: "test-token",
				TokenType:   "bearer",
				ExpiresIn:   3600,
			},
			expectedErr: false,
		},
		{
			name: "Empty Access Token",
			serverResponse: map[string]interface{}{
				"access_token": "",
			},
			serverStatus: http.StatusOK,
			expectedAuth: Auth{},
			expectedErr:  true,
		},
		{
			name:           "Invalid JSON Response",
			serverResponse: nil,
			serverStatus:   http.StatusOK,
			expectedAuth:   Auth{},
			expectedErr:    true,
		},
		{
			name:           "Non 200 Status Code",
			serverResponse: map[string]interface{}{"error": "unauthorized"},
			serverStatus:   http.StatusUnauthorized,
			expectedAuth:   Auth{},
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
					t.Errorf("expected contnet type application/x-www-form-urlencoded, got %s", ct)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != nil {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					w.Write([]byte("invalid json"))
				}
			}))
			defer server.Close()

			envs := map[string]string{
				"CLIENT_ID":     "test-client-id",
				"CLIENT_SECRET": "test-client-secret",
				"GRANT_TYPE":    "client-credentials",
			}

			for k, v := range envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			client := &http.Client{}
			auth, err := PostForAuth(client, server.URL)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if auth != tt.expectedAuth {
				t.Errorf("expected auth: %v, got: %v", tt.expectedAuth, auth)
			}
		})
	}
}
