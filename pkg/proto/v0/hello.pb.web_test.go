package proto_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var invalidEndpoints = []string{"api", "api/v0", "greet", ""}

func TestPOSTCorrectEndpointVariousData(t *testing.T) {
	type TestStruct struct {
		testDataName         string
		requestBody          string
		expectedStatusCode   int
		responseBody         string
		expectedResponseBody string
	}

	var tests = []TestStruct{
		{"no-data", `{}`, 400, "", "missing a name\n"},
		{"wrong key", `{"naam": "Milan"}`, 412, "", `unknown field "naam" in proto.GreetRequest` + "\n"},
		{"empty body", ``, 412, "", "EOF\n"},
		{"invalid json", `{"name":"Milan"{}`, 412, "", "invalid character '{' after object key:value pair\n"},
		{"data is int", `{"name":23}`, 412, "", "json: cannot unmarshal number into Go value of type string\n"},
		{"data is json", `{"name":{age: 23}}`, 412, "", "invalid character 'a' looking for beginning of object key string\n"},
		{"additional data", `{"name":"Milan", "surname":"Bahadur"}`, 412, "", `unknown field "surname" in proto.GreetRequest` + "\n"},
		{"value missing", `{"name":""}`, 400, "", "missing a name\n"},
		{"ASCII name", `{"name":"Milan"}`, 201, "", `{"message":"Hello Milan"}` + "\n"},
		{"UTF name", `{"name":"मिलन"}`, 201, "", `{"message":"Hello मिलन"}` + "\n"},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {

			rr := sendRequest(t, "POST", "/api/v0/greet", testCase.requestBody)
			assertResult(t, rr, testCase.expectedStatusCode, testCase.expectedResponseBody)
		})
	}
}

func TestPOSTIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			rr := sendRequest(t, "POST", endpoint, `{"name":"Milan"}`)
			assertResult(t, rr, 404, "404 page not found\n")
		})
	}
}

func TestGETIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			rr := sendRequest(t, "GET", endpoint, "")
			assertResult(t, rr, 404, "404 page not found\n")
		})
	}
}

func sendRequest(t *testing.T, method, endpoint, data string) *httptest.ResponseRecorder {
	var reader = strings.NewReader(data)
	req, err := http.NewRequest(method, endpoint, reader)
	assert.Nil(t, err)

	r := chi.NewRouter()
	proto.RegisterHelloWeb(r, svc.NewService())

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertResult(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int, expectedBody string) {
	assert.Equal(t, expectedBody, rr.Body.String(), "response body not as expected")
	assert.Equal(t, expectedStatusCode, rr.Code, "response code not as expected")
}
