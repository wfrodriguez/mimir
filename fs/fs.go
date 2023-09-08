package fs

import (
	"os"
	"path/filepath"
)

func FindFileByExt(path, ext string) ([]string, error) {
	var files = make([]string, 0)
	if ext[0] != '.' {
		ext = "." + ext
	}

	err := filepath.Walk(path, func(p string, fi os.FileInfo, er error) error {
		if er != nil {
			return er
		}

		if !fi.IsDir() && filepath.Ext(p) == ext {
			files = append(files, p)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
