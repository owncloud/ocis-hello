package proto

import (
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

var host = "http://localhost:9105/"

func TestMain(m *testing.M) {
	if os.Getenv("SYSTEM_UNDER_TEST") != "" {
		host = os.Getenv("SYSTEM_UNDER_TEST")
	}
	code := m.Run()

	os.Exit(code)
}

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
		{"value missing", `{"name":""}`, 400, "", "missing a name\n"},
		{"ASCII name", `{"name":"Milan"}`, 201, "", `{"message":"Hello Milan"}` + "\n"},
		{"UTF name", `{"name":"मिलन"}`, 201, "", `{"message":"Hello मिलन"}` + "\n"},
	}

	endpoint := "api/v0/greet"
	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			res := sendRequest(t, "POST", endpoint, testCase.requestBody)
			assertResult(t, res, testCase.expectedStatusCode, testCase.expectedResponseBody)
		})
	}
}

func TestPOSTIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			res := sendRequest(t, "POST", endpoint, `{"name":"Milan"}`)
			assertResult(t, res, 404, "404 page not found\n")
		})
	}
}

func TestGETIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			res := sendRequest(t, "GET", endpoint, "")
			assertResult(t, res, 404, "404 page not found\n")
		})
	}
}

func sendRequest(t *testing.T, method, endpoint, data string) *http.Response {
	var reader = strings.NewReader(data)
	request, err := http.NewRequest(method, host+endpoint, reader)
	if err != nil {
		t.Error(err)
	}
	
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	return res
}

func assertResult(t *testing.T, res *http.Response, expectedStatusCode int, expectedBody string) {
	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, string(body), expectedBody, "response body not as expected")
	assert.Equal(t, res.StatusCode, expectedStatusCode, "response code not as expected")
}
