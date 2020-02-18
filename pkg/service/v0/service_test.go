package svc

import (
	"context"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"gotest.tools/assert"
	"testing"
)

func TestVariousCorrectData(t *testing.T) {
	type TestStruct struct {
		testDataName           string
		requestString          string
		expectedResponseString string
		expectedError		   string
	}

	var tests = []TestStruct{
		{"no-data",  "", "", "missing a name"},
		{"ASCII name", "Milan", "Hello Milan", ""},
		{"UTF name", "मिलन",`Hello मिलन`, ""},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			ctx := context.Background()
			svc := NewService()
			req := v0proto.GreetRequest{Name: testCase.requestString}
			res := v0proto.GreetResponse{}
			err := svc.Greet(ctx, &req, &res)
			if err == nil && testCase.expectedError != "" {
				t.Error("Expected error '" + testCase.expectedError + "' but no error returned")

			} else if err != nil || testCase.expectedError != "" {
				assert.Equal(t, err.Error(), testCase.expectedError)
			}
			assert.Equal(t, res.Message, testCase.expectedResponseString)
		})
	}
}
