package file_utils

import (
	"bytes"
	"io"
)

const DefaultChunkSize = int64(1024)

type BackwardScanner struct {
	r   io.ReaderAt
	pos int64
	err error
	buf []byte
}

func NewBackwardScanner(r io.ReaderAt, pos int64) *BackwardScanner {
	return &BackwardScanner{r: r, pos: pos}
}

func (s *BackwardScanner) readMore() {
	if s.pos == 0 {
		s.err = io.EOF
		return
	}
	size := DefaultChunkSize
	if size > s.pos {
		size = s.pos
	}
	s.pos -= size
	buf2 := make([]byte, size, size+int64(len(s.buf)))

	// ReadAt attempts to read full buff!
	_, s.err = s.r.ReadAt(buf2, int64(s.pos))
	if s.err == nil {
		s.buf = append(buf2, s.buf...)
	}
}

func (s *BackwardScanner) Line() (line string, start int64, err error) {
	if s.err != nil {
		return "", 0, s.err
	}
	for {
		lineStart := bytes.LastIndexByte(s.buf, '\n')
		if lineStart >= 0 { // we have a complete line
			var line string
			line, s.buf = string(dropCR(s.buf[lineStart+1:])), s.buf[:lineStart]
			return line, s.pos + int64(lineStart) + int64(1), nil
		}
		// need more data
		s.readMore()
		if s.err != nil {
			if s.err == io.EOF {
				if len(s.buf) > 0 {
					return string(dropCR(s.buf)), 0, nil
				}
			}
			return "", 0, s.err
		}
	}
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
