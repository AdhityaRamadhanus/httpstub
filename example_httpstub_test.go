package httpstub_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AdhityaRamadhanus/httpstub"
)

func Example_default() {
	srv := httpstub.NewStubServer()
	srv.StubRequest(http.MethodGet, "/healthz")
	defer srv.Close()

	url := fmt.Sprintf("%s%s", srv.URL(), "/healthz")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Body: %s\n", string(respBody))
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode)

	// Output:
	// Response Body: OK
	// Response Status Code: 200
}

func Example_responseJSON() {
	srv := httpstub.NewStubServer()
	responseJSON := map[string]interface{}{
		"success": true,
	}
	srv.StubRequest(http.MethodGet, "/healthz", httpstub.WithResponseBodyJSON(responseJSON))
	defer srv.Close()

	url := fmt.Sprintf("%s%s", srv.URL(), "/healthz")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Body: %s\n", string(respBody))
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode)
	// Output:
	// Response Body: {"success":true}
	// Response Status Code: 200
}

func Example_responseJSONFile() {
	srv := httpstub.NewStubServer()
	responseJSONPath := "example_response.json"
	srv.StubRequest(http.MethodGet, "/healthz", httpstub.WithResponseBodyJSONFile(responseJSONPath))
	defer srv.Close()

	url := fmt.Sprintf("%s%s", srv.URL(), "/healthz")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Body: %s\n", string(respBody))
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode)
	// Output:
	// Response Body: {"success":true}
	// Response Status Code: 200
}
