package v0alpha1svc

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrMissingName = errors.New("missing a name")
)

// Basic returns a naive implementation of service.
func Basic() Service {
	return basic{}
}

type basic struct{}

func (b basic) Greet(_ context.Context, name string) (string, error) {
	if name == "" {
		return "", ErrMissingName
	}

	return fmt.Sprintf("Hello %s", name), nil
}
