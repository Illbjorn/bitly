package settings

import (
	"bytes"
	"io"
	"strconv"

	"github.com/illbjorn/bitly/internal/safe"
)

var buf = bytes.NewBuffer(make([]byte, 0, 1024))

func List() string {
	const (
		colMaxW = 8
	)

	buf.Reset()

	// Boilerplate Header
	buf.WriteString("SETTINGS\n")
	const prefix = "  "
	alignedRow(buf, colMaxW, prefix, "Name", "Value", "Description")
	alignedRow(buf, colMaxW, prefix, "----", "-----", "-----------")

	// String Settings
	for _, s := range settings.strings {
		alignedRow(buf, colMaxW, prefix, s.name, s.value, s.description)
	}

	// Int Settings
	for _, i := range settings.ints {
		var itoa = strconv.Itoa
		alignedRow(buf, colMaxW, prefix, i.name, itoa(i.value), i.description)
	}

	// Bool Settings
	for _, b := range settings.bools {
		var btoa = strconv.FormatBool
		alignedRow(buf, colMaxW, prefix, b.name, btoa(b.value), b.description)
	}

	return safe.Btos(buf.Bytes())
}

func alignedRow(w io.Writer, align int, prefix string, cols ...string) {
	const space = " "
	const newline = "\n"
	w.Write(safe.Stob(prefix))
	for _, col := range cols {
		w.Write(safe.Stob(col))
		for range max(0, align-len(col)) {
			w.Write(safe.Stob(space))
		}
	}
	w.Write(safe.Stob(newline))
}
