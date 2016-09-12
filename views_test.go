package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test helpers

func NewRequest(method, endpoint, body string) (*httptest.ResponseRecorder, *http.Request) {
	responseRecorder := httptest.NewRecorder()
	bodyReader := strings.NewReader(body)
	request := httptest.NewRequest(method, endpoint, bodyReader)
	return responseRecorder, request
}

func httpStatusCodeTest(rw *httptest.ResponseRecorder, expectedStatusCode int) func(*testing.T) {
	return func(t *testing.T) {
		if rw.Code != expectedStatusCode {
			t.Errorf("Expected status code `%i`, actual `%i`", expectedStatusCode, rw.Code)
		}
	}
}

func httpBodyTest(rw *httptest.ResponseRecorder, expectedBody string) func(*testing.T) {
	return func(t *testing.T) {
		actualBody := rw.Body.String()
		if actualBody != expectedBody {
			t.Errorf("Expected response: `%s`\nActual response: `%s`", expectedBody, actualBody)
		}
	}
}

func checkResponse(t *testing.T, rw *httptest.ResponseRecorder, statusCode int, body string) {
	if statusCode != 0 {
		t.Run("StatusCode", httpStatusCodeTest(rw, statusCode))
	}

	if body != "" {
		t.Run("ResponseBody", httpBodyTest(rw, body))
	}
}

// Actual tests

func TestHealthcheckResponse(t *testing.T) {
	router := NewMockRouter()
	rw, request := NewRequest("GET", "/healthcheck", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Status":"ok"}`)
}

func TestAddURL(t *testing.T) {
	router := NewMockRouter()

	rw, request := NewRequest("POST", "/create", `{"Url": "http://www.nationalreview.com"}`)
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Url":"bs1I92"}`)
}

func TestFetchURL(t *testing.T) {
	router := NewMockRouter()

	rw, request := NewRequest("GET", "/foobar", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 301, "")

	rw, request = NewRequest("GET", "/redsox", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 404, "")
}

func TestUrlStats(t *testing.T) {
	router := NewMockRouter()

	rw, request := NewRequest("GET", "/stats/ghjk", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Count":387,"Days":{"2015-07-22T00:00:00Z":14,"2015-11-03T00:00:00Z":76,"2016-01-03T00:00:00Z":31}}`)

	rw, request = NewRequest("GET", "/stats/china", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 404, "")
}
