package lib

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Attempts to copy the file from the given path
// to the given path.
// This method will create the destination folder
// if it does not yet exist and then copy the
// file if the destination file does not yet exist.
// If any step in the mean time fails then the error
// is logged and renaming the file is stopped - however
// the execution of the program is not aborted.
func AttemptCopy(before string, after string) bool {
	log.Printf("[---] %-20v -> %v", before, after)

	// Ensure that the directory exists.
	dir := filepath.Dir(after)
	if !pathExists(dir) {
		log.Printf("[---] Creating directory %v", dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Printf("[ERR] Cannot created directory %v: %v", dir, err)
			return false
		}
	}

	// Abort if destination file already exists
	if pathExists(after) {
		log.Printf("[ERR] %-20v -> %v: Destination file already exists", before, after)
		return false
	}

	// Rename the file
	// TODO Use something else than shelling out to cp
	err := exec.Command("cp", "-rf", before, after).Run()
	if err != nil {
		log.Printf("[ERR] %v -> %v Could not rename file: %v", before, after, err)
	}
	return true
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
