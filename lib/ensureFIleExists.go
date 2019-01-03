package lib

import (
	"os"
)

func EnsureFileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		emptyFile, err := os.Create(file)
		emptyFile.Close()
		return err
	}
	return nil
}
