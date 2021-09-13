package file_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kroppt/voxels/modules/file"
)

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestModuleNew(t *testing.T) {
	t.Run("return is non-nil", func(t *testing.T) {
		mod := file.New()
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestModuleGetReadCloser(t *testing.T) {

	t.Run("open empty file", func(t *testing.T) {
		expectErr := io.EOF
		expectN := 0
		mod := file.New()

		readCloser, err := mod.GetReadCloser("test-empty.txt")

		if err != nil {
			t.Fatal(err)
		}
		if readCloser == nil {
			t.Fatalf("expected non-nil reader")
		}
		defer readCloser.Close()
		buf := make([]byte, 8)
		n, err := readCloser.Read(buf)
		if !errors.Is(err, expectErr) {
			t.Fatalf("expected %q but got %q", io.EOF, err)
		}
		if n != expectN {
			t.Fatalf("expected %v bytes read but got %v bytes read", expectN, n)
		}
	})

	t.Run("open single-line file", func(t *testing.T) {
		expectN := 6
		mod := file.New()

		readCloser, err := mod.GetReadCloser("test-single-line.txt")

		if err != nil {
			t.Fatal(err)
		}
		if readCloser == nil {
			t.Fatalf("expected non-nil reader")
		}
		defer readCloser.Close()
		buf := make([]byte, 16)
		n, err := readCloser.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if n != expectN {
			t.Fatalf("expected %v bytes read but got %v bytes read", expectN, n)
		}
	})

	t.Run("open multi-line file", func(t *testing.T) {
		expect := []string{"Line 1.", "123abc", "This is a test.", ""}
		mod := file.New()

		readCloser, err := mod.GetReadCloser("test-multi-line.txt")

		if err != nil {
			t.Fatal(err)
		}
		if readCloser == nil {
			t.Fatalf("expected non-nil reader")
		}
		defer readCloser.Close()
		buf := make([]byte, 64)
		n, err := readCloser.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		content := string(buf[:n])
		content = strings.ReplaceAll(content, "\r\n", "\n")
		lines := strings.Split(content, "\n")
		if !stringSliceEqual(expect, lines) {
			t.Fatalf("expected %v but got %v", expect, lines)
		}
	})

}
