package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUtils(t *testing.T) {
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

	t.Run("WriteError", func(t *testing.T) {
		runTest := func(status int, err error) {
			w := httptest.NewRecorder()
			WriteError(w, status, err)

			if status := w.Result().StatusCode; status != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if response["error"] != err.Error() {
				t.Errorf("expected error message %q, got %q", err.Error(), response["error"])
			}
		}

		runTest(http.StatusBadRequest, fmt.Errorf("test error"))
	})

	t.Run("WriteJSONWithCookie", func(t *testing.T) {
		runTest := func(status int, v any, cookies []*http.Cookie) {
			w := httptest.NewRecorder()

			WriteJSONWithCookie(w, status, v, cookies)

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

	t.Run("WriteJSON", func(t *testing.T) {
		runTest := func(status int, v any) {
			w := httptest.NewRecorder()

			WriteJSON(w, status, v)

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

