package httpstub

import (
	"net/http"
)

type matcher struct {
	Method         string
	Path           string
	RequestHeaders map[string]string
}

func (m *matcher) matchRequest(req *http.Request) bool {
	url := req.URL.String()
	if m.Path != url {
		return false
	}

	if m.Method != req.Method {
		return false
	}

	if m.RequestHeaders != nil {
		isSubset := true
		for header, value := range m.RequestHeaders {
			headerVal := req.Header.Get(header)
			if headerVal != value {
				isSubset = false
				break
			}
		}
		if !isSubset {
			return false
		}
	}

	return true
}
