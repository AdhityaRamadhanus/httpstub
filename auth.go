package httpstub

import (
	"encoding/base64"
	"fmt"
)

func WithBasicAuth(username, password string) func(spec *spec) {
	credentials := fmt.Sprintf("%s:%s", username, password)
	credentials = base64.StdEncoding.EncodeToString([]byte(credentials))
	return func(spec *spec) {
		spec.RequestHeaders["Authorization"] = credentials
	}
}

func WithBearerToken(token string) func(spec *spec) {
	return func(spec *spec) {
		spec.RequestHeaders["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
}
