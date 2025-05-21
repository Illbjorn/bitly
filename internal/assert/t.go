package assert

import "testing"

var t *testing.T

func SetT(_t *testing.T) {
	t = _t
}
