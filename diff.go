// Package diff checks whether files and io.Readers have identical contents.
package diff

import (
	"bytes"
	"io"
	"os"
)

const chunkSize = 4096

// Readers compares the contents of two io.Readers.
// The return value of different is true if and only if there are no errors
// in reading r1 and r2 (io.EOF excluded) and r1 and r2 are
// byte-for-byte identical.
func Readers(r1, r2 io.Reader) (different bool, err error) {
	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)
	for {
		short1 := false
		n1, err := io.ReadFull(r1, buf1)
		switch err {
		case io.EOF, io.ErrUnexpectedEOF:
			short1 = true
		case nil:
		default:
			return true, err
		}
		short2 := false
		n2, err := io.ReadFull(r2, buf2)
		switch err {
		case io.EOF, io.ErrUnexpectedEOF:
			short2 = true
		case nil:
		default:
			return true, err
		}
		if short1 != short2 || n1 != n2 {
			return true, nil
		}
		if !bytes.Equal(buf1[:n1], buf2[:n1]) {
			return true, nil
		}
		if short1 {
			return false, nil
		}
	}
}

// Files compares the contents of file1 and file2.
// Files first compares file length before looking at the contents.
func Files(file1, file2 string) (different bool, err error) {
	f1, err := os.Open(file1)
	if err != nil {
		return true, err
	}
	defer f1.Close()
	f2, err := os.Open(file2)
	if err != nil {
		return true, err
	}
	defer f2.Close()

	// Compare the size of the files.
	n1, err := f1.Seek(0, io.SeekEnd)
	if err != nil {
		return true, err
	}
	n2, err := f2.Seek(0, io.SeekEnd)
	if err != nil {
		return true, err
	}
	if n1 != n2 {
		return true, nil
	}
	if _, err := f1.Seek(0, io.SeekStart); err != nil {
		return true, err
	}
	if _, err := f2.Seek(0, io.SeekStart); err != nil {
		return true, err
	}

	// Otherwise compare the contents.
	return Readers(f1, f2)
}
