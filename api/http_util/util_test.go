package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apiError "github.com/beka-birhanu/finance-go/api/error"
)

// TestUtils contains tests for the utility functions in the httputil package.
func TestUtils(t *testing.T) {
	// Test ParseJSON function
	t.Run("ParseJSON", func(t *testing.T) {
		runTestCase := func(input []byte, expected string, shouldErr bool) {
			r := &http.Request{
				Body: io.NopCloser(bytes.NewReader(input)),
			}

			var result struct {
				Name string `json:"name"`
			}

			err := ParseJSON(r, &result)
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

	// Test RespondError function
	t.Run("RespondError", func(t *testing.T) {
		runTest := func(testErr apiError.Error) {
			w := httptest.NewRecorder()
			RespondError(w, testErr)

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

		testErr := apiError.NewBadRequest("test error")
		runTest(testErr)
	})

	// Test RespondWithCookies function
	t.Run("RespondWithCookies", func(t *testing.T) {
		runTest := func(status int, v any, cookies []*http.Cookie) {
			w := httptest.NewRecorder()

			RespondWithCookies(w, status, v, cookies)

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

	// Test Respond function
	t.Run("Respond", func(t *testing.T) {
		runTest := func(status int, v any) {
			w := httptest.NewRecorder()

			Respond(w, status, v)

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
}

