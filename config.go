package httpstub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Config is a helper to create spec that StubServer knows how to handle
type Config func(spec *spec)

//WithRequestHeaders set request headers that will be matched on an incoming matched request to StubServer
func WithRequestHeaders(headers map[string]string) Config {
	return func(spec *spec) {
		spec.requestMatcher.RequestHeaders = headers
	}
}

//WithResponseHeaders set response headers that will be sent to client on an incoming matched request to StubServer
func WithResponseHeaders(headers map[string]string) Config {
	return func(spec *spec) {
		spec.ResponseHeaders = headers
	}
}

//WithResponseCode set response status code that will be sent to client on an incoming matched request to StubServer
func WithResponseCode(statusCode int) Config {
	return func(spec *spec) {
		spec.ResponseCode = statusCode
	}
}

//WithResponseBody set response body(bytes) that will be sent to client on an incoming matched request to StubServer
func WithResponseBody(body []byte) Config {
	return func(spec *spec) {
		spec.ResponseBody = body
	}
}

//WithResponseBodyString set response body(string) that will be sent to client on an incoming matched request to StubServer
func WithResponseBodyString(body string) Config {
	return func(spec *spec) {
		spec.setResponseHeader("Content-Type", "text/plain; charset=utf-8")
		spec.ResponseBody = []byte(body)
	}
}

//WithResponseBodyJSON set response body(map[string]interface{}) that will be sent to client on an incoming matched request to StubServer
func WithResponseBodyJSON(body map[string]interface{}) Config {
	return func(spec *spec) {
		if bodyBytes, err := json.Marshal(body); err == nil {
			spec.setResponseHeader("Content-Type", "application/json; charset=utf-8")
			spec.ResponseBody = bodyBytes
		}
	}
}

//WithResponseBodyJSONFile set response body(path to json) that will be sent to client on an incoming matched request to StubServer
func WithResponseBodyJSONFile(path string) Config {
	return func(spec *spec) {
		body, err := ioutil.ReadFile(path)
		if err == nil {
			spec.setResponseHeader("Content-Type", "application/json; charset=utf-8")
			spec.ResponseBody = body
		}
	}
}

//WithResponseBodyFile set response body(path to file) that will be sent to client on an incoming matched request to StubServer
func WithResponseBodyFile(path string) Config {
	return func(spec *spec) {
		body, err := ioutil.ReadFile(path)
		if err == nil {
			spec.ResponseBody = body
		}
	}
}

//WithBasicAuth will match only request with provided basic auth credentials
func WithBasicAuth(username, password string) Config {
	credentials := fmt.Sprintf("%s:%s", username, password)
	credentials = base64.StdEncoding.EncodeToString([]byte(credentials))
	return func(spec *spec) {
		spec.setRequestHeader("Authorization", fmt.Sprintf("Basic %s", credentials))
	}
}

//WithBearerToken will match only request with provided bearer token
func WithBearerToken(token string) Config {
	return func(spec *spec) {
		spec.setRequestHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}
