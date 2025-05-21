package assert

func Equal[T comparable](v1, v2 T) {
	t.Helper()
	if v1 != v2 {
		t.Logf("%v does not equal %v", v1, v2)
		t.Fail()
	}
}

func True(cond bool) {
	t.Helper()
	if !cond {
		t.Logf("expected true value")
		t.Fail()
	}
}

func False(cond bool) {
	t.Helper()
	if cond {
		t.Logf("expected false value")
		t.Fail()
	}
}

func Error(err error) {
	t.Helper()
	if err == nil {
		t.Logf("expected error, found none")
		t.Fail()
	}
}

func NoError(err error) {
	t.Helper()
	if err != nil {
		t.Logf("expected no error, found: %s", err)
		t.Fail()
	}
}
