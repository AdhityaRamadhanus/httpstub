//Package httpstub is wrapper package of httptest Server to make it easier to stub http request to a third-party
// when testing a module
package httpstub

import "net/http"

// Spec is a configuration on h
type spec struct {
	// Request (matcher)
	requestMatcher matcher

	// Response
	ResponseBody    []byte
	ResponseHeaders map[string]string
	ResponseCode    int
}

var (
	defaultOKSpec = spec{
		ResponseBody: []byte("OK"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: http.StatusOK,
	}

	notAllowedSpec = spec{
		ResponseBody: []byte("Method not allowed"),
		ResponseHeaders: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
		ResponseCode: http.StatusMethodNotAllowed,
	}
)

func (s *spec) matchRequest(req *http.Request) bool {
	return s.requestMatcher.matchRequest(req)
}

func (s *spec) setRequestHeader(header, value string) {
	if s.requestMatcher.RequestHeaders == nil {
		s.requestMatcher.RequestHeaders = map[string]string{}
	}
	s.requestMatcher.RequestHeaders[header] = value
}

func (s *spec) setResponseHeader(header, value string) {
	if s.ResponseHeaders == nil {
		s.ResponseHeaders = map[string]string{}
	}
	s.ResponseHeaders[header] = value
}
