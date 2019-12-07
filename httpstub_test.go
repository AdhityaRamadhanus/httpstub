package httpstub_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/AdhityaRamadhanus/httpstub"
)

type checkFunc func(t *testing.T, res *http.Response, body string)
type routeSpec struct {
	Method  string
	Path    string
	Options []httpstub.Config
}

func newServer(t *testing.T, routeSpecs []routeSpec) (*httpstub.StubServer, func()) {
	t.Helper()
	srv := httpstub.NewStubServer()
	for _, spec := range routeSpecs {
		srv.StubRequest(
			spec.Method,
			spec.Path,
			spec.Options...,
		)
	}

	teardown := func() {
		srv.Close()
	}

	return srv, teardown
}

var (
	checkHTTPResponseHeader = func(key, value string) checkFunc {
		return func(t *testing.T, res *http.Response, body string) {
			if res.Header.Get(key) != value {
				t.Errorf("responseHeader[%q] = %s; want %s", key, res.Header.Get(key), value)
			}
		}
	}

	checkHTTPResponseBody = func(val string) checkFunc {
		return func(t *testing.T, res *http.Response, body string) {
			if body != val {
				t.Errorf("responseBody = %s; want %s", body, val)
			}
		}
	}

	checkHTTPResponseCode = func(statusCode int) checkFunc {
		return func(t *testing.T, res *http.Response, body string) {
			if res.StatusCode != statusCode {
				t.Errorf("responseCode = %d; want %d", res.StatusCode, statusCode)
			}
		}
	}
)

func TestStubRequest_matcher(t *testing.T) {
	// preparation
	routeSpecs := []routeSpec{
		{http.MethodGet, "/healthz/get", nil},
		{http.MethodPost, "/healthz/post", nil},
		{http.MethodHead, "/healthz/head", nil},
		{http.MethodGet, "/healthz/custom-headers", []httpstub.Config{
			httpstub.WithRequestHeaders(map[string]string{"Authorization": "Test"})},
		},
		{http.MethodPost, "/healthz/basic-auth", []httpstub.Config{
			httpstub.WithBasicAuth("test", "test")},
		},
		{http.MethodGet, "/healthz/conflicted-headers", []httpstub.Config{
			httpstub.WithBasicAuth("test", "test"),
			httpstub.WithBearerToken("bearer token")},
		},
	}

	server, teardown := newServer(t, routeSpecs)
	defer teardown()

	testCases := []struct {
		Method         string
		Path           string
		RequestHeaders map[string]string

		// used in assertion
		checkFuncs []checkFunc
	}{
		{
			Method: http.MethodGet,
			Path:   "/healthz/get",
			checkFuncs: []checkFunc{
				checkHTTPResponseBody("OK"),
				checkHTTPResponseCode(http.StatusOK),
			},
		},
		{
			Method: http.MethodPost,
			Path:   "/healthz/post",
			checkFuncs: []checkFunc{
				checkHTTPResponseBody("OK"),
				checkHTTPResponseCode(http.StatusOK),
			},
		},
		{
			Method: http.MethodHead,
			Path:   "/healthz/head",
			checkFuncs: []checkFunc{
				checkHTTPResponseBody(""),
				checkHTTPResponseCode(http.StatusOK),
			},
		},
		{
			Method:         http.MethodGet,
			Path:           "/healthz/custom-headers",
			RequestHeaders: map[string]string{"Authorization": "Test"},
			checkFuncs: []checkFunc{
				checkHTTPResponseBody("OK"),
				checkHTTPResponseCode(http.StatusOK),
			},
		},
		{
			Method:         http.MethodPost,
			Path:           "/healthz/basic-auth",
			RequestHeaders: map[string]string{"Authorization": "Basic dGVzdDp0ZXN0"},
			checkFuncs: []checkFunc{
				checkHTTPResponseBody("OK"),
				checkHTTPResponseCode(http.StatusOK),
			},
		},
	}

	t.Run("parallel_group", func(t *testing.T) {
		for _, testCase := range testCases {
			testCase := testCase
			t.Run(testCase.Path, func(t *testing.T) {
				t.Parallel()
				url := fmt.Sprintf("%s%s", server.URL(), testCase.Path)
				req, err := http.NewRequest(testCase.Method, url, nil)
				if err != nil {
					t.Fatalf("http.NewRequest() err = %s; want nil", err)
				}
				for header, value := range testCase.RequestHeaders {
					req.Header.Set(header, value)
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("http.DefaultClient.Do(req) err = %s; want nil", err)
				}
				defer resp.Body.Close()

				respBodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("ioutil.ReadALl(resp.Body) err = %s; want nil", err)
				}

				for _, checkFunc := range testCase.checkFuncs {
					checkFunc(t, resp, string(respBodyBytes))
				}
			})
		}
	})
}
