package writer

import "io"

type NopWriter struct {
	io.Writer
}

func (*NopWriter) Write(buf []byte) (int, error) {
	return len(buf), nil
}
