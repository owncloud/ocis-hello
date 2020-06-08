package proto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
)
type TestStruct struct {
	testDataName string
	name         string
	expected     string
}

func TestRequestString(t *testing.T) {
	var tests = []TestStruct{
		{"ASCII", "Milan", `name:"Milan"`},
		{"UTF", "मिलन", `name:"मिलन"`},
		{"empty", "" , ``},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			request := proto.GreetRequest{Name: testCase.name}
			assert.Equal(t, testCase.expected, request.String())
		})
	}
}

func TestResponseString(t *testing.T) {
	var tests = []TestStruct{
		{"ASCII", "Milan", `message:"Milan"`},
		{"UTF", "मिलन", `message:"मिलन"`},
		{"empty", "" , ``},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			response := proto.GreetResponse{Message: testCase.name}
			assert.Equal(t, testCase.expected, response.String())
		})
	}
}
