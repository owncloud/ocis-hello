package proto_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/stretchr/testify/assert"

	svc "github.com/owncloud/ocis-hello/pkg/service/v0"
)

var service = grpc.Service{}

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("hello"),
		grpc.Address("localhost:9992"),
	)

	err := proto.RegisterHelloHandler(service.Server(), svc.NewService())
	if err != nil {
		log.Fatalf("could not register HelloHandler: %v", err)
	}
	service.Server().Start()
}

type ErrorMessage struct {
	Id     string
	Code   int
	Detail string
	Status string
}

func TestCorrectService(t *testing.T) {
	type TestStruct struct {
		testDataName     string
		name             string
		expectedError    ErrorMessage
		expectedResponse interface{}
	}

	var tests = []TestStruct{
		{
			"value missing",
			"",
			ErrorMessage{
				"go.micro.client",
				500,
				"missing a name",
				"Internal Server Error",
			},
			nil,
		},
		{"ASCII name",
			"Milan",
			ErrorMessage{},
			&proto.GreetResponse{
				Message: "Hello Milan",
			},
		},
		{"special characters",
			`$&/\#.`,
			ErrorMessage{},
			&proto.GreetResponse{
				Message: `Hello $&/\#.`,
			},
		},
		{"UTF name",
			"मिलन",
			ErrorMessage{},
			&proto.GreetResponse{
				Message: "Hello मिलन",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			request := proto.GreetRequest{Name: testCase.name}
			client := service.Client()
			cl := proto.NewHelloService("com.owncloud.api.hello", client)
			response, err := cl.Greet(context.Background(), &request)
			if err != nil || (ErrorMessage{}) != testCase.expectedError {
				assert.Nil(t, response)
				assert.Error(t, err)
				var errorData ErrorMessage
				json.Unmarshal([]byte(err.Error()), &errorData)
				assert.Equal(t, testCase.expectedError.Id, errorData.Id)
				assert.Equal(t, testCase.expectedError.Code, errorData.Code)
				assert.Equal(t, testCase.expectedError.Detail, errorData.Detail)
				assert.Equal(t, testCase.expectedError.Status, errorData.Status)
			}
			if testCase.expectedResponse != nil {
				assert.Nil(t, err)
				// TODO: this should test the whole response, not just the message. BUT: response body doesn't get closed
				// when only checking the `response` itself and there is no way to close it, other than accessing it as of now.
				assert.Equal(t, testCase.expectedResponse.(*proto.GreetResponse).Message, response.Message)
			}
		})
	}
}

func TestWrongService(t *testing.T) {
	var tests = []string{
		"com.owncloud.api",
		"com.owncloud.api.greet",
		"com.owncloud.hello",
		`com/owncloud/api/hello`,
		"",
	}

	for _, testCase := range tests {
		t.Run(testCase, func(t *testing.T) {
			request := proto.GreetRequest{Name: "Milan"}
			client := service.Client()
			cl := proto.NewHelloService(testCase, client)
			response, err := cl.Greet(context.Background(), &request)
			assert.Nil(t, response)
			assert.Error(t, err)
			var errorData ErrorMessage
			json.Unmarshal([]byte(err.Error()), &errorData)
			assert.Equal(t, 500, errorData.Code)
			assert.Equal(t, "service "+testCase+": not found", errorData.Detail)
			assert.Equal(t, "Internal Server Error", errorData.Status)
		})
	}
}
