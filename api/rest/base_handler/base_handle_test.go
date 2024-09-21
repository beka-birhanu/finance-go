package baseapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestBaseHandler(t *testing.T) {
	handler := &BaseHandler{}

	// Test ValidatedBody method
	t.Run("ValidatedBody", func(t *testing.T) {
		runTestCase := func(input []byte, shouldErr bool) {
			r := &http.Request{
				Body: io.NopCloser(bytes.NewReader(input)),
			}

			var result struct {
				Name string `json:"name" validate:"required"`
			}

			err := handler.ValidatedBody(r, &result)
			if (err != nil) != shouldErr {
				t.Fatalf("expected error: %v, got: %v", shouldErr, err)
			}
		}

		cases := []struct {
			input     []byte
			shouldErr bool
		}{
			{
				input:     []byte(`{"name":"John"}`),
				shouldErr: false,
			},
			{
				input:     []byte(``),
				shouldErr: true,
			},
			{
				input:     []byte(`{"name":""}`),
				shouldErr: true,
			},
			{
				input:     []byte(`not-a-json`),
				shouldErr: true,
			},
		}

		for _, c := range cases {
			runTestCase(c.input, c.shouldErr)
		}
	})

	// Test Problem method
	t.Run("Problem", func(t *testing.T) {
		runTest := func(testErr errapi.Error) {
			w := httptest.NewRecorder()
			handler.Problem(w, testErr)

			if status := w.Result().StatusCode; status != testErr.StatusCode() {
				t.Errorf("expected status code %d, got %d", testErr.StatusCode(), status)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if response["error"] != testErr.Error() {
				t.Errorf("expected error message %q, got %q", testErr.Error(), response["error"])
			}
		}

		testErr := errapi.NewBadRequest("test error")
		runTest(testErr)
	})

	// Test ParseJSON method
	t.Run("ParseJSON", func(t *testing.T) {
		runTestCase := func(input []byte, expected string, shouldErr bool) {
			r := &http.Request{
				Body: io.NopCloser(bytes.NewReader(input)),
			}

			var result struct {
				Name string `json:"name"`
			}

			err := handler.ParseJSON(r, &result)
			if (err != nil) != shouldErr {
				t.Fatalf("expected error: %v, got: %v", shouldErr, err)
			}
			if !shouldErr && result.Name != expected {
				t.Errorf("expected name: %s, got: %s", expected, result.Name)
			}
		}

		cases := []struct {
			input     []byte
			expected  string
			shouldErr bool
		}{
			{
				input:     []byte(`{"name":"John"}`),
				expected:  "John",
				shouldErr: false,
			},
			{
				input:     []byte(``),
				expected:  "",
				shouldErr: true,
			},
			{
				input:     []byte(`not-a-json`),
				expected:  "",
				shouldErr: true,
			},
		}

		for _, c := range cases {
			runTestCase(c.input, c.expected, c.shouldErr)
		}
	})

	// Test RespondError method
	t.Run("RespondError", func(t *testing.T) {
		runTest := func(testErr errapi.Error) {
			w := httptest.NewRecorder()
			handler.RespondError(w, testErr)

			if status := w.Result().StatusCode; status != testErr.StatusCode() {
				t.Errorf("expected status code %d, got %d", testErr.StatusCode(), status)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if response["error"] != testErr.Error() {
				t.Errorf("expected error message %q, got %q", testErr.Error(), response["error"])
			}
		}

		testErr := errapi.NewBadRequest("test error")
		runTest(testErr)
	})

	// Test RespondWithCookies method
	t.Run("RespondWithCookies", func(t *testing.T) {
		runTest := func(status int, v any, cookies []*http.Cookie) {
			w := httptest.NewRecorder()

			handler.RespondWithCookies(w, status, v, cookies)

			if status := w.Result().StatusCode; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if response["key"] != "value" {
				t.Errorf("expected response key 'key' to be 'value', got %q", response["key"])
			}

			cookiesReceived := w.Result().Cookies()
			if len(cookiesReceived) != 1 || cookiesReceived[0].Name != "testCookie" || cookiesReceived[0].Value != "testValue" {
				t.Errorf("expected cookie 'testCookie' with value 'testValue', got %+v", cookiesReceived)
			}
		}

		cookies := []*http.Cookie{
			{Name: "testCookie", Value: "testValue"},
		}
		runTest(http.StatusOK, map[string]string{"key": "value"}, cookies)
	})

	// Test Respond method
	t.Run("Respond", func(t *testing.T) {
		runTest := func(status int, v any) {
			w := httptest.NewRecorder()

			handler.Respond(w, status, v)

			if status := w.Result().StatusCode; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			contentType := w.Result().Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected content type 'application/json', got %q", contentType)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if response["key"] != "value" {
				t.Errorf("expected response key 'key' to be 'value', got %q", response["key"])
			}
		}

		runTest(http.StatusOK, map[string]string{"key": "value"})
	})

	// Test BaseURL method
	t.Run("BaseURL", func(t *testing.T) {
		r := &http.Request{
			Host: "example.com",
		}

		baseURL := handler.BaseURL(r)
		expected := "http://example.com"
		if baseURL != expected {
			t.Errorf("expected baseURL: %s, got: %s", expected, baseURL)
		}

		r.TLS = &tls.ConnectionState{}
		baseURL = handler.BaseURL(r)
		expected = "https://example.com"
		if baseURL != expected {
			t.Errorf("expected baseURL: %s, got: %s", expected, baseURL)
		}
	})

	// Test UUIDParam method
	t.Run("UUIDParam", func(t *testing.T) {
		validUUID := uuid.New().String()
		invalidUUID := "invalid-uuid"

		runTestCase := func(paramName, paramValue string, shouldErr bool) {
			r := mux.SetURLVars(&http.Request{}, map[string]string{
				"id": paramValue,
			})

			_, err := handler.UUIDParam(r, paramName)
			if (err != nil) != shouldErr {
				t.Fatalf("expected error: %v, got: %v at %s param name", shouldErr, err, paramName)
			}
		}

		runTestCase("id", validUUID, false)
		runTestCase("id", invalidUUID, true)
		runTestCase("missing", validUUID, true)
	})
}
