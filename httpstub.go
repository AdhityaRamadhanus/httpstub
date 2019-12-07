package httpstub

import (
	"net/http"
	"net/http/httptest"
)

//StubServer is a server that stub request coming and response according to a spec (see config)
type StubServer struct {
	httpServer *httptest.Server
	specs      []spec
}

//NewStubServer Create and Run StubServer
func NewStubServer() *StubServer {
	stubServer := &StubServer{
		specs: []spec{},
	}
	stubServer.httpServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var foundSpec = notAllowedSpec
		for _, spec := range stubServer.specs {
			if spec.matchRequest(req) {
				foundSpec = spec
			}
		}

		stubServer.createResponse(res, foundSpec)
	}))

	return stubServer
}

func (s *StubServer) createResponse(res http.ResponseWriter, spec spec) {
	responseHeaders := defaultOKSpec.ResponseHeaders
	// merge
	for header, value := range spec.ResponseHeaders {
		responseHeaders[header] = value
	}
	for header, value := range responseHeaders {
		res.Header().Set(header, value)
	}

	responseCode := defaultOKSpec.ResponseCode
	if spec.ResponseCode != 0 {
		responseCode = spec.ResponseCode
	}
	res.WriteHeader(responseCode)

	responseBody := defaultOKSpec.ResponseBody
	if spec.ResponseBody != nil {
		responseBody = spec.ResponseBody
	}

	// HEAD METHOD HANDLED BY GO, SO IT WILL NOT RETURN ANYTHING

	res.Write(responseBody)
}

//StubRequest takes method, path and configs to create a spec that will be matched on server
func (s *StubServer) StubRequest(method, path string, options ...Config) {
	spec := spec{
		requestMatcher: matcher{
			Method: method,
			Path:   path,
		},
	}

	for _, option := range options {
		option(&spec)
	}

	s.specs = append(s.specs, spec)
}

//URL return StubServer URL, your client should make a request to this URL
func (s *StubServer) URL() string {
	if s.httpServer == nil {
		return "unknown"
	}

	return s.httpServer.URL
}

//Close the StubServer
func (s *StubServer) Close() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}
