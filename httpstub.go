package httpstub

import (
	"net/http"
	"net/http/httptest"
)

type Config func(spec *Spec)

var (
	defaultSpec = Spec{
		ResponseBody: []byte("OK"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: http.StatusOK,
	}

	notAllowedSpec = Spec{
		ResponseBody: []byte("Method not allowed"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: http.StatusMethodNotAllowed,
	}
)

type Server struct {
	server *httptest.Server
	specs  []Spec
}

func WithRequestHeaders(headers map[string]string) Config {
	return func(spec *Spec) {
		spec.RequestHeaders = headers
	}
}

func defaultHTTPHandler(res http.ResponseWriter, req *http.Request, spec Spec) {
	responseHeaders := defaultSpec.ResponseHeaders
	// merge
	for header, value := range spec.ResponseHeaders {
		responseHeaders[header] = value
	}
	for header, value := range responseHeaders {
		res.Header().Set(header, value)
	}

	responseCode := defaultSpec.ResponseCode
	if spec.ResponseCode != 0 {
		responseCode = spec.ResponseCode
	}
	res.WriteHeader(responseCode)

	responseBody := defaultSpec.ResponseBody
	if spec.ResponseBody != nil {
		responseBody = spec.ResponseBody
	}

	res.Write(responseBody)
}

func (s *Server) StubRequest(method, path string, options ...Config) {
	spec := Spec{
		Method: method,
		Path:   path,
	}

	for _, option := range options {
		option(&spec)
	}

	s.specs = append(s.specs, spec)
}

func (s *Server) URL() string {
	if s.server != nil {
		return s.server.URL
	}

	return "unknown"
}

func (s *Server) Start() {
	s.server = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		url := req.URL.String()
		var foundSpec = notAllowedSpec
		for _, spec := range s.specs {
			if spec.Path != url {
				continue
			}

			if spec.Method != req.Method {
				continue
			}

			if spec.RequestHeaders != nil {
				isSubset := true
				for header, value := range spec.RequestHeaders {
					headerVal := req.Header.Get(header)
					if headerVal != value {
						isSubset = false
						break
					}
				}
				if !isSubset {
					continue
				}
			}

			foundSpec = spec
		}

		defaultHTTPHandler(res, req, foundSpec)
	}))
}

func (s *Server) Close() {
	if s.server != nil {
		s.server.Close()
	}
}
