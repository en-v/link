package link

import (
	"os"
	"path/filepath"
)

func FilexExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LocDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
