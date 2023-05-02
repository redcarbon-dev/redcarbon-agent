package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

const dirPermission = 0o755

func Untar(archive io.Reader, target string) error {
	gzr, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}

	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if header == nil {
			continue
		}

		target := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, dirPermission); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			dirPath := filepath.Dir(target)
			if _, err := os.Stat(dirPath); err != nil {
				if err := os.MkdirAll(dirPath, dirPermission); err != nil {
					return err
				}
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}
}
