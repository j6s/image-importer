package lib

import (
	"bufio"
	"os"
)

func ReadLines(file string) (error, []string) {
	existing := make([]string, 0)
	filePointer, err := os.Open(file)
	if err != nil {
		return err, existing
	}
	defer filePointer.Close()

	scanner := bufio.NewScanner(filePointer)
	for scanner.Scan() {
		existing = append(existing, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err, existing
	}

	return nil, existing
}
