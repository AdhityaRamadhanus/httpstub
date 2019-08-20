package httpstub

import (
	"encoding/json"
	"io/ioutil"
)

func WithResponseHeaders(headers map[string]string) Config {
	return func(spec *Spec) {
		spec.ResponseHeaders = headers
	}
}

func WithResponseCode(statusCode int) Config {
	return func(spec *Spec) {
		spec.ResponseCode = statusCode
	}
}

func WithResponseBody(body []byte) Config {
	return func(spec *Spec) {
		spec.ResponseBody = body
	}
}

func WithResponseBodyString(body string) Config {
	return func(spec *Spec) {
		spec.SetResponseHeader("Content-Type", "text/plain; charset=utf-8")
		spec.ResponseBody = []byte(body)
	}
}

func WithResponseBodyJSON(body map[string]interface{}) Config {
	return func(spec *Spec) {
		var bodyBytes []byte
		err := json.Unmarshal(bodyBytes, body)
		if err == nil {
			spec.SetResponseHeader("Content-Type", "application/json; charset=utf-8")
			spec.ResponseBody = bodyBytes
		}
	}
}

func WithResponseBodyJSONFile(path string) Config {
	return func(spec *Spec) {
		body, err := ioutil.ReadFile(path)
		if err == nil {
			spec.SetResponseHeader("Content-Type", "application/json; charset=utf-8")
			spec.ResponseBody = body
		}
	}
}

func WithResponseBodyFile(path string) Config {
	return func(spec *Spec) {
		body, err := ioutil.ReadFile(path)
		if err == nil {
			spec.ResponseBody = body
		}
	}
}
