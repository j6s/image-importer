package main

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"strings"

	lib "github.com/j6s/image-importer/lib"

	"github.com/mostlygeek/go-exiftool"
)

type FileProperties struct {
	Date      string
	Time      string
	Year      string
	Month     string
	Day       string
	Hour      string
	Minute    string
	Second    string
	Name      string
	Extension string
}

func main() {
	settings := parseSettings()
	fmt.Printf("Importing images using %d workers\n", settings.Workers)

	imported := 0
	lib.WorkerPool(
		func(file string) {
			if doCopyFile(file, &settings) {
				imported = imported + 1
			}
		},
		func(channel chan string) {
			for _, file := range settings.Files {
				channel <- file
			}
		},
		settings.Workers,
	)

	lib.WriteLinesToFile(settings.BlacklistPath, settings.Blacklist)
	log.Printf("Imported %d images", imported)
}

func doCopyFile(before string, settings *Settings) bool {

	// Read EXIF
	exif, err := exiftool.Extract(before)
	if err != nil {
		log.Printf("[ERR] Could not parse exif for file %v : %v", before, err.Error())
		return false
	}

	// Build new filename
	created, _ := exif.CreateDate()
	extension := path.Ext(before)
	var afterFilename bytes.Buffer
	settings.FilenameTemplate.Execute(&afterFilename, FileProperties{
		Date:      created.Format("2006-01-02"),
		Time:      created.Format("15-04-05"),
		Year:      created.Format("2006"),
		Month:     created.Format("01"),
		Day:       created.Format("02"),
		Hour:      created.Format("15"),
		Minute:    created.Format("04"),
		Second:    created.Format("05"),
		Name:      strings.Replace(path.Base(before), extension, "", 1),
		Extension: extension[1:],
	})
	after := fmt.Sprintf("%s/%s", settings.To, afterFilename.String())

	// Copy file
	copied := lib.AttemptCopy(before, after)
	if copied {
		settings.Blacklist = append(settings.Blacklist, path.Base(before))
	}
	return copied
}
