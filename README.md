# httpstub
A wrapper package for an easier way to stub http request in golang

<p>
  <a href="#Usage">Usage |</a>
  <a href="#licenses">License</a>
  <br><br>
</p>

Usage
-----
* You will need access token to use the api
* generate access token
```go
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
```

License
----

GPL © [Adhitya Ramadhanus]

