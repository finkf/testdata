package testdata // import "github.com/finkf/testdata"

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var update *bool

func init() {
	update = flag.Bool("update", false, "update gold file contents")

}

var errHandle = func(err error) { panic(err) }

// SetErrHandle updates the error handling function.
// The error handling function is called by Must
// only iff the err is not nil.
func SetErrHandle(h func(err error)) {
	errHandle = h
}

// Must calls the ErrHandle if the given error is not nil.
func Must(err error) {
	if err != nil {
		errHandle(err)
	}
}

// Dir defines the name of the testdata directory.
const Dir = "testdata"

// File returns the path to a test file in the
// testdata directory.
func File(name string) string {
	return filepath.Join(Dir, name)
}

// Bytes reads the given test file into a byte slice.
func Bytes(name string) []byte {
	is, err := os.Open(File(name))
	Must(err)
	defer is.Close()
	bs, err := ioutil.ReadAll(is)
	Must(err)
	return bs
}

// String reads the given test file into a string.
func String(name string) string {
	return string(Bytes(name))
}

// Reader reads the given test file and returns a Reader.
func Reader(name string) io.Reader {
	return bytes.NewBuffer(Bytes(name))
}

// maybeUpdate checks if the given gold file should be updated.
// It calls the given function to update.
func maybeUpdate(name string, f func(string)) {
	flag.Parse()
	if *update {
		f(name)
	}
}

// GoldString checks if the given content and the content
// of the gold file are the same. It calls t.Fatalf iff
// they are not the same.
// If go test was called with the -update flag, the contents
// of the gold file are updated with the given content.
func GoldString(t *testing.T, got, name string) {
	t.Helper()
	maybeUpdate(name, func(name string) {
		UpdateString(got, name)
	})
	TestGoldString(t, got, name)
}

// TestGoldString compares the content of a gold file
// with the given content. It calls t.Fatalf if the strings
// are not the same.
func TestGoldString(t *testing.T, got, name string) {
	t.Helper()
	want := String(name)
	if got == want {
		return
	}
	t.Fatalf(`gold file: %q
expected:
%q
got:
%q
`, File(name), want, got)
}

// UpdateString updates the condent of a gold file.
func UpdateString(content, name string) {
	UpdateBytes([]byte(content), name)
}

// UpdateBytes updates the condent of a gold file.
func UpdateBytes(content []byte, name string) {
	err := ioutil.WriteFile(File(name), content, 0666)
	Must(err)
}

// UpdateReader updates the condent of a gold file.
func UpdateReader(content io.Reader, name string) {
	bs, err := ioutil.ReadAll(content)
	Must(err)
	UpdateBytes(bs, name)
}
