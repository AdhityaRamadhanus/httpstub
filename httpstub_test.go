package httpstub_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/AdhityaRamadhanus/httpstub"
)

func TestStubRequest(t *testing.T) {
	// preparation
	stubSpecs := []struct {
		Method  string
		Path    string
		Options []httpstub.Config
	}{
		{
			Method: http.MethodGet,
			Path:   "/healthz",
		},
		{
			Method: http.MethodGet,
			Path:   "/healthz",
			Options: []httpstub.Config{
				httpstub.WithRequestHeaders(map[string]string{"Authorization": "Test"}),
				httpstub.WithResponseBodyString("OK With Authorization"),
			},
		},
		{
			Method: http.MethodPost,
			Path:   "/healthz/status",
			Options: []httpstub.Config{
				httpstub.WithBasicAuth("test", "test"),
				httpstub.WithResponseBodyString("OK With Authorization"),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/healthz",
			Options: []httpstub.Config{
				httpstub.WithRequestHeaders(map[string]string{"Authorization": "Test"}),
				httpstub.WithResponseBodyString("OK With Authorization"),
			},
		},
	}

	testCases := []struct {
		Name           string
		Method         string
		Path           string
		RequestHeaders map[string]string

		// used in assertion
		ResponseHeaders map[string]string
		ResponseCode    int
		ResponseBody    string
	}{
		{
			Name:         "Return 200",
			Method:       http.MethodGet,
			Path:         "/healthz",
			ResponseCode: http.StatusOK,
			ResponseBody: "OK",
		},
		{
			Name:         "Return 405",
			Method:       http.MethodGet,
			Path:         "/healthzzzzz",
			ResponseCode: http.StatusMethodNotAllowed,
			ResponseBody: "Method not allowed",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			srv := httpstub.NewStubServer()
			for _, spec := range stubSpecs {
				srv.StubRequest(
					spec.Method,
					spec.Path,
					spec.Options...,
				)
			}
			defer srv.Close()

			url := fmt.Sprintf("%s%s", srv.URL(), testCase.Path)
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

			wantBody := testCase.ResponseBody
			gotBody := string(respBodyBytes)
			if gotBody != wantBody {
				t.Errorf("StubRequest(spec) body = %s; want %s", gotBody, wantBody)
			}

			wantCode := testCase.ResponseCode
			gotCode := resp.StatusCode
			if gotCode != wantCode {
				t.Errorf("StubRequest(spec) statusCode = %d; want %d", gotCode, wantCode)
			}
		})
	}
}
