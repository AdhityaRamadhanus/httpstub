package httpstub

import (
	"encoding/base64"
	"fmt"
)

func WithBasicAuth(username, password string) Config {
	credentials := fmt.Sprintf("%s:%s", username, password)
	credentials = base64.StdEncoding.EncodeToString([]byte(credentials))
	return func(spec *Spec) {
		if spec.RequestHeaders == nil {
			spec.RequestHeaders = map[string]string{}
		}
		spec.RequestHeaders["Authorization"] = credentials
	}
}

func WithBearerToken(token string) Config {
	return func(spec *Spec) {
		if spec.RequestHeaders == nil {
			spec.RequestHeaders = map[string]string{}
		}
		spec.RequestHeaders["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
}
