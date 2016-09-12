package main

import (
	"github.com/gocraft/web"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test helpers

func NewRequest(method, endpoint, body string, pathParams map[string]string) (*httptest.ResponseRecorder, *web.Request) {
	responseRecorder := httptest.NewRecorder()
	bodyReader := strings.NewReader(body)
	if pathParams == nil {
		pathParams = make(map[string]string)
	}
	httpRequest := httptest.NewRequest(method, endpoint, bodyReader)
	webRequest := &web.Request{Request: httpRequest, PathParams: pathParams}
	return responseRecorder, webRequest
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
	server := CreateMockServer()
	rw, request := NewRequest("GET", "/healthcheck", "", nil)
	server.healthcheck(rw, request)
	checkResponse(t, rw, 200, `{"Status":"ok"}`)
}

func TestAddURL(t *testing.T) {
	server := CreateMockServer()

	rw, request := NewRequest("POST", "/create", `{"Url": "http://www.nationalreview.com"}`, nil)
	server.addUrl(rw, request)
	checkResponse(t, rw, 200, `{"Url":"bs1I92"}`)
}

func TestFetchURL(t *testing.T) {
	server := CreateMockServer()

	rw, request := NewRequest("GET", "/foobar", "", map[string]string{"path": "foobar"})
	server.fetchUrl(rw, request)
	checkResponse(t, rw, 301, "")

	rw, request = NewRequest("GET", "/redsox", "", map[string]string{"path": "redsox"})
	server.fetchUrl(rw, request)
	checkResponse(t, rw, 404, "")
}

func TestUrlStats(t *testing.T) {
	server := CreateMockServer()

	rw, request := NewRequest("GET", "/stats/ghjk", "", map[string]string{"path": "ghjk"})
	server.urlStats(rw, request)
	checkResponse(t, rw, 200, `{"Count":387,"Days":{"2015-07-22T00:00:00Z":14,"2015-11-03T00:00:00Z":76,"2016-01-03T00:00:00Z":31}}`)

	rw, request = NewRequest("GET", "/stats/china", "", map[string]string{"path": "china"})
	server.urlStats(rw, request)
	checkResponse(t, rw, 404, "")
}
