package cpio

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func safeJoin(dst, name string) (string, error) {
	name = filepath.ToSlash(name)
	name = strings.TrimPrefix(name, "/")
	clean := filepath.Clean(name)
	if clean == "." || clean == "" {
		return "", errors.New("cpio: invalid name")
	}
	if strings.HasPrefix(clean, "..") || strings.Contains(clean, "../") {
		return "", errors.New("cpio: path traversal")
	}
	return filepath.Join(dst, filepath.FromSlash(clean)), nil
}

func permsFromMode(mode uint32) fs.FileMode {
	return fs.FileMode(mode & 0777)
}

// UnpackToDir unpacks a CPIO archive into dst.
func UnpackToDir(r io.Reader, dst string) error {
	rd := NewReader(r)
	for {
		e, err := rd.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		outPath, err := safeJoin(dst, e.Name)
		if err != nil {
			return err
		}

		isDir := (e.Mode & 0170000) == 0040000
		if isDir {
			if err := os.MkdirAll(outPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, e.Data, permsFromMode(e.Mode)); err != nil {
			return err
		}
	}
}
