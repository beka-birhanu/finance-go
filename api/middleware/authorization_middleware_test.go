package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/dgrijalva/jwt-go"
)

type MockJwtService struct {
	DecodeTokenFunc func(token string) (jwt.MapClaims, error)
}

func (m *MockJwtService) Generate(user *usermodel.User) (string, error) {
	return "", nil
}

func (m *MockJwtService) Decode(token string) (jwt.MapClaims, error) {
	return m.DecodeTokenFunc(token)
}

func TestAuthorizationMiddleware(t *testing.T) {
	tests := []struct {
		name                 string
		setCookie            bool
		mockDecodeTokenFunc  func(token string) (jwt.MapClaims, error)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Missing access token",
			setCookie:            false,
			mockDecodeTokenFunc:  nil,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "Authorization token required\n",
		},
		{
			name:      "Invalid access token",
			setCookie: true,
			mockDecodeTokenFunc: func(token string) (jwt.MapClaims, error) {
				return nil, errors.New("invalid token")
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "Invalid token\n",
		},
		{
			name:      "Valid access token",
			setCookie: true,
			mockDecodeTokenFunc: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"user_id": "123"}, nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "Hello, authorized user!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJwtService := &MockJwtService{
				DecodeTokenFunc: tt.mockDecodeTokenFunc,
			}

			mw := Authorization(mockJwtService, true)

			// Create a handler to be wrapped by the middleware
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims := r.Context().Value(ContextUserClaims)
				if claims != nil {
					if _, err := w.Write([]byte("Hello, authorized user!\n")); err != nil {
						t.Error("error in writing ")
					}
				} else {
					if _, err := w.Write([]byte("Hello, guest!\n")); err != nil {
						t.Error("error in writing ")
					}
				}
			})

			// Wrap the handler with the middleware
			wrappedHandler := mw(handler)

			// Create a new HTTP request
			req := httptest.NewRequest("GET", "/", nil)

			// Set the access token cookie if required
			if tt.setCookie {
				req.AddCookie(&http.Cookie{Name: "accessToken", Value: "dummyToken"})
			}

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Serve the request
			wrappedHandler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tt.expectedResponseBody {
				t.Errorf("expected body %q, got %q", tt.expectedResponseBody, rr.Body.String())
			}
		})
	}
}
