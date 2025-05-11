package suite

import "testing"

type Suite struct {
	*testing.T
}

func New(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	return &Suite{
		T: t,
	}
}
