package suite

import "testing"

type Suite struct {
	*testing.T
}

func New() *Suite {
	return &Suite{}
}
