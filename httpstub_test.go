package httpstub_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/AdhityaRamadhanus/httpstub"
	"github.com/stretchr/testify/assert"
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
			Name:         "Simple spec, return 200",
			Method:       http.MethodGet,
			Path:         "/healthz",
			ResponseCode: http.StatusOK,
			ResponseBody: "OK",
		},
		{
			Name:         "Path and method not found in spec, return 405",
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
			srv := httpstub.Server{}
			srv.Start()
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
			for header, value := range testCase.RequestHeaders {
				req.Header.Set(header, value)
			}

			client := http.Client{}
			resp, err := client.Do(req)

			// assert
			if assert.NoError(t, err, "Should not return error") {
				defer resp.Body.Close()
				respBodyBytes, _ := ioutil.ReadAll(resp.Body)
				assert.Equal(t, testCase.ResponseBody, string(respBodyBytes))
				assert.Equal(t, testCase.ResponseCode, resp.StatusCode)
			}
		})
	}
}
