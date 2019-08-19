package httpstub_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/AdhityaRamadhanus/httpstub"
	"github.com/stretchr/testify/assert"
)

func TestStubRequest(t *testing.T) {
	srv := httpstub.Server{}
	srv.StubRequest(
		http.MethodGet,
		"/healthz",
		httpstub.WithRequestHeaders(map[string]string{
			"Authorization": "Test",
		}),
		httpstub.WithResponseHeaders(map[string]string{
			"X-HttpStub": "Test",
		}),
		httpstub.WithResponseBody(strings.NewReader("OK")),
		httpstub.WithResponseCode(http.StatusOK),
	)
	srv.Start()
	defer srv.Close()

	url := fmt.Sprintf("%s/healthz", srv.URL())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Test")
	client := http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Should not return error")
	defer resp.Body.Close()
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "Test", resp.Header.Get("X-HttpStub"))
	assert.Equal(t, "OK", string(respBodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
