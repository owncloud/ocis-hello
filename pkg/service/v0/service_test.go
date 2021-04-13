package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
)

func TestHello_Greet(t *testing.T) {
	tests := []struct {
		name                 string
		req                  string
		expectedErrorMessage interface{}
	}{
		{"simple", "simple", nil},
		{"UTF", "मिलन", nil},
		{"special char", `%&# /\`, nil},
		{"empty", "", "missing a name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			req := &v0proto.GreetRequest{Name: tt.req}
			var rsp = &v0proto.GreetResponse{}

			err := s.Greet(context.Background(), req, rsp)

			if tt.expectedErrorMessage != nil || err != nil {
				assert.EqualError(t, err, tt.expectedErrorMessage.(string))
			} else {
				assert.Equal(t, "Hello "+tt.req, rsp.Message)
			}
		})
	}
}
