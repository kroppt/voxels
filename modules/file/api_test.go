package file_test

import (
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
	t.Parallel()
	t.Run("return is non-nil", func(t *testing.T) {
		mod := file.New()
		if mod == nil {
			t.Fatal("expected non-nil return")
		}
	})
}

func TestModuleGetReadCloser(t *testing.T) {
	t.Parallel()

	t.Run("open empty file", func(t *testing.T) {
		t.Parallel()
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
		buf, err := io.ReadAll(readCloser)
		if err != nil {
			t.Fatal(err)
		}
		n := len(buf)
		if n != expectN {
			t.Fatalf("expected %v bytes read but got %v bytes read", expectN, n)
		}
	})

	t.Run("open single-line file", func(t *testing.T) {
		t.Parallel()
		expect := "abc123"
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
		buf, err := io.ReadAll(readCloser)
		if err != nil {
			t.Fatal(err)
		}
		content := string(buf)
		n := len(buf)
		if content != expect {
			t.Fatalf("expected %v but got %v", expect, content)
		}
		if n != expectN {
			t.Fatalf("expected %v bytes read but got %v bytes read", expectN, n)
		}
	})

	t.Run("open multi-line file", func(t *testing.T) {
		t.Parallel()
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
		buf, err := io.ReadAll(readCloser)
		if err != nil {
			t.Fatal(err)
		}
		content := strings.ReplaceAll(string(buf), "\r\n", "\n")
		lines := strings.Split(content, "\n")
		if !stringSliceEqual(expect, lines) {
			t.Fatalf("expected %v but got %v", expect, lines)
		}
	})

}
