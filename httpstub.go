package httpstub

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

type spec struct {
	// Request (matcher)
	Method         string
	Path           string
	RequestHeaders map[string]string

	// Response
	ResponseBody    io.Reader
	ResponseHeaders map[string]string
	ResponseCode    int
}

var (
	defaultSpec = spec{
		ResponseBody: strings.NewReader("OK"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: 200,
	}

	notAllowedSpec = spec{
		ResponseBody: strings.NewReader("Method not allowed"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: 405,
	}
)

type Server struct {
	server *httptest.Server
	specs  []spec
}

func WithRequestHeaders(headers map[string]string) func(spec *spec) {
	return func(spec *spec) {
		spec.RequestHeaders = headers
	}
}

func WithResponseHeaders(headers map[string]string) func(spec *spec) {
	return func(spec *spec) {
		spec.ResponseHeaders = headers
	}
}

func WithResponseCode(statusCode int) func(spec *spec) {
	return func(spec *spec) {
		spec.ResponseCode = statusCode
	}
}

func WithResponseBody(body io.Reader) func(spec *spec) {
	return func(spec *spec) {
		spec.ResponseBody = body
	}
}

func WithResponseBodyFromFile(path string) func(spec *spec) {
	return func(spec *spec) {
		body, err := os.Open(path)
		if err != nil {
			spec.ResponseBody = body
		}
	}
}

func defaultHTTPHandler(res http.ResponseWriter, req *http.Request, spec spec) {
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
	responseBytes, _ := ioutil.ReadAll(responseBody)
	res.Write(responseBytes)
}

func (s *Server) StubRequest(method, path string, options ...func(spec *spec)) {
	spec := spec{
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

			defaultHTTPHandler(res, req, spec)
			return
		}

		// not allowed
		defaultHTTPHandler(res, req, notAllowedSpec)
	}))
}

func (s *Server) Close() {
	if s.server != nil {
		s.server.Close()
	}
}
