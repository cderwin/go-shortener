package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
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
	_, router := NewMockRouter()
	rw, request := NewRequest("GET", "/healthcheck", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Status":"ok"}`)
}

func TestAddURL(t *testing.T) {
	_, router := NewMockRouter()

	rw, request := NewRequest("POST", "/create", `{"Url": "http://www.nationalreview.com"}`)
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Url":"bs1I92"}`)
}

func TestFetchURL(t *testing.T) {
	server, router := NewMockRouter()

	// Test handler redirects as expected
	rw, request := NewRequest("GET", "/foobar", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 301, "")

	// Test 404 returned when url does not exist
	rw, request = NewRequest("GET", "/redsox", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 404, "")

	// Test hits are incremented when URL is hit
	shortURL, _ := server.Redis.SaveURL("https://news.ycombinator.com")
	rw, request = NewRequest("GET", "/"+shortURL, "")
	router.ServeHTTP(rw, request)
	t.Run("checkIncremented", func(t *testing.T) {
		expected := Hits{Count: 1, Days: map[time.Time]int{MockNow: 1}}
		actual, _ := server.Redis.GetHits(shortURL)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected: %+v\nActual: %+v\n", expected, actual)
		}
	})
}

func TestUrlStats(t *testing.T) {
	_, router := NewMockRouter()

	rw, request := NewRequest("GET", "/stats/ghjk", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 200, `{"Count":387,"Days":{"2015-07-22T00:00:00Z":14,"2015-11-03T00:00:00Z":76,"2016-01-03T00:00:00Z":31}}`)

	rw, request = NewRequest("GET", "/stats/china", "")
	router.ServeHTTP(rw, request)
	checkResponse(t, rw, 404, "")
}
