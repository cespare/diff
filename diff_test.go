package diff

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReaders(t *testing.T) {
	for _, tc := range [][2]string{
		{"", ""},
		{"a", "a"},
		{"a", ""},
		{"a", "aa"},
		{"a", "b"},
	} {
		got, err := Readers(strings.NewReader(tc[0]), strings.NewReader(tc[1]))
		if err != nil {
			t.Fatal(err)
		}
		want := tc[0] == tc[1]
		if got != want {
			t.Fatalf("comparing %q and %q; got %t; want %t", tc[0], tc[1], got, want)
		}
	}
}

func TestFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "diff-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	var b1 []byte
	b2 := make([]byte, chunkSize)
	b3 := make([]byte, chunkSize)
	b3[chunkSize-1] = 123
	b4 := make([]byte, chunkSize-1)
	b5 := make([]byte, chunkSize+1)
	writeFile(t, dir, b1, "f1")
	writeFile(t, dir, b2, "f2")
	writeFile(t, dir, b3, "f3")
	writeFile(t, dir, b4, "f4")
	writeFile(t, dir, b5, "f5")

	for _, tc := range []struct {
		f1, f2    string
		identical bool
	}{
		{"f1", "f1", true},
		{"f1", "f2", false},
		{"f2", "f2", true},
		{"f2", "f3", false},
		{"f2", "f4", false},
		{"f2", "f5", false},
	} {
		got, err := Files(filepath.Join(dir, tc.f1), filepath.Join(dir, tc.f2))
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.identical {
			t.Fatalf("comparing %s and %s; got %t; want %t", tc.f1, tc.f2, got, tc.identical)
		}
	}
}

func writeFile(t *testing.T, dir string, b []byte, name string) {
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := f.Write(b); err != nil {
		t.Fatal(err)
	}
}
