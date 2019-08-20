package httpstub

type Spec struct {
	// Request (matcher)
	Method         string
	Path           string
	RequestHeaders map[string]string

	// Response
	ResponseBody    []byte
	ResponseHeaders map[string]string
	ResponseCode    int
}

func (s *Spec) SetRequestHeader(header, value string) {
	if s.RequestHeaders == nil {
		s.RequestHeaders = map[string]string{}
	}
	s.RequestHeaders[header] = value
}

func (s *Spec) SetResponseHeader(header, value string) {
	if s.ResponseHeaders == nil {
		s.ResponseHeaders = map[string]string{}
	}
	s.ResponseHeaders[header] = value
}
