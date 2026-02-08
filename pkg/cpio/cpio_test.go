package cpio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPackAndReadRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "bin", "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := PackDir(tmp, &buf, WithMTimeUnix(0)); err != nil {
		t.Fatal(err)
	}

	rd := NewReader(bytes.NewReader(buf.Bytes()))
	seen := map[string][]byte{}
	for {
		e, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		seen[e.Name] = e.Data
	}

	if got := string(seen["bin/hello.txt"]); got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestUnpackToDir(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "etc"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "etc", "conf"), []byte("x=y"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := PackDir(tmp, &buf, WithMTimeUnix(0)); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	if err := UnpackToDir(bytes.NewReader(buf.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "etc", "conf"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "x=y" {
		t.Fatalf("got %q", string(b))
	}
}
