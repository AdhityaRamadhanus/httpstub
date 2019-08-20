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
srv := httpstub.Server{}
srv.StubRequest(
  http.MethodGet,
  "/healthz",
  httpstub.WithResponseHeaders(map[string]string{
    "Content-Type": "application/json; charset=utf-8",
  }),
  httpstub.WithResponseBodyFile("../test/json/healthz.json"),
)
srv.Start()
defer srv.Close()

url := fmt.Sprintf("%s%s", srv.URL(), "/healthz")
req, err := http.NewRequest(testCase.Method, url, nil)

client := http.Client{}
resp, err := client.Do(req)
// do something with resp
```

License
----

GPL © [Adhitya Ramadhanus]

