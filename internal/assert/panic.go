package assert

func Panics(fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r == nil {
			t.Logf("expected panic")
			t.Fail()
		}
	}()
	fn()
}
