package lib

import (
	"fmt"
	"os"
)

func WriteLinesToFile(path string, lines []string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	for _, line := range lines {
		_, err := file.Write([]byte(fmt.Sprintf("%s\n", line)))
		if err != nil {
			return err
		}
	}

	return nil
}
