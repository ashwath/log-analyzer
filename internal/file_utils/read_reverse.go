package file_utils

import (
	"bytes"
	"io"
)

const DefaultChunkSize = 1024

type Scanner struct {
	r   io.ReaderAt
	pos int
	err error
	buf []byte
}

func NewScanner(r io.ReaderAt, pos int) *Scanner {
	return &Scanner{r: r, pos: pos}
}

func (s *Scanner) readMore() {
	if s.pos == 0 {
		s.err = io.EOF
		return
	}
	size := DefaultChunkSize
	if size > s.pos {
		size = s.pos
	}
	s.pos -= size
	buf2 := make([]byte, size, size+len(s.buf))

	// ReadAt attempts to read full buff!
	_, s.err = s.r.ReadAt(buf2, int64(s.pos))
	if s.err == nil {
		s.buf = append(buf2, s.buf...)
	}
}

func (s *Scanner) Line() (line string, start int, err error) {
	if s.err != nil {
		return "", 0, s.err
	}
	for {
		lineStart := bytes.LastIndexByte(s.buf, '\n')
		if lineStart >= 0 { // we have a complete line
			var line string
			line, s.buf = string(s.buf[lineStart+1:]), s.buf[:lineStart]
			return line, s.pos + lineStart + 1, nil
		}
		// need more data
		s.readMore()
		if s.err != nil {
			if s.err == io.EOF {
				if len(s.buf) > 0 {
					return string(s.buf), 0, nil
				}
			}
			return "", 0, s.err
		}
	}
}
