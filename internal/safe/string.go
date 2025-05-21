package safe

import "unsafe"

func Btos(b []byte) (s string) {
	if len(b) == 0 {
		return ""
	}
	sd := unsafe.SliceData(b)
	return unsafe.String(sd, len(b))
}

func Stob(s string) (b []byte) {
	sd := unsafe.StringData(s)
	return unsafe.Slice(sd, len(s))
}
