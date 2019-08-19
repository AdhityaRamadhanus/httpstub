package httpstub

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Spec struct {
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
	defaultSpec = Spec{
		ResponseBody: strings.NewReader("OK"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: 200,
	}

	notAllowedSpec = Spec{
		ResponseBody: strings.NewReader("Method not allowed"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: 405,
	}
)

type Server struct {
	server *httptest.Server
	Specs  []Spec
}

func WithRequestHeaders(headers map[string]string) func(spec *Spec) {
	return func(spec *Spec) {
		spec.RequestHeaders = headers
	}
}

func WithResponseHeaders(headers map[string]string) func(spec *Spec) {
	return func(spec *Spec) {
		spec.ResponseHeaders = headers
	}
}

func WithResponseCode(statusCode int) func(spec *Spec) {
	return func(spec *Spec) {
		spec.ResponseCode = statusCode
	}
}

func WithResponseBody(body io.Reader) func(spec *Spec) {
	return func(spec *Spec) {
		spec.ResponseBody = body
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
	responseBytes, _ := ioutil.ReadAll(responseBody)
	res.Write(responseBytes)
}

func (s *Server) StubRequest(method, path string, options ...func(spec *Spec)) {
	spec := Spec{
		Method: method,
		Path:   path,
	}

	for _, option := range options {
		option(&spec)
	}

	s.Specs = append(s.Specs, spec)
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
		for _, spec := range s.Specs {
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
