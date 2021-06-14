package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHello_Greet(t *testing.T) {
	tests := []struct {
		name string
		req  string
	}{
		{"simple", "simple"},
		{"UTF", "मिलन"},
		{"special char", `%&# /\`},
		{"empty", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BasicGreeter{phraseSource: StaticPhraseSource{DefaultPhrase}}

			greeting := s.Greet("", tt.req)

			assert.Equal(t, "Hello "+tt.req, greeting)
		})
	}
}
