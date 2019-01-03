package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	lib "github.com/j6s/image-importer/lib"
)

type Settings struct {
	Workers          int
	From             string
	To               string
	BlacklistPath    string
	Blacklist        []string
	Files            []string
	FilesRaw         []string
	FilenameTemplate template.Template
}

var (
	_from      = flag.String("from", "", "The source directory of the files to import. All files in this directory (matching the given extension) will be processed.")
	_to        = flag.String("to", "", "The destination directory of the files. The files will be copied to this folder.")
	_index     = flag.String("index", ".imported", "Name of the index file in the source folder.")
	_extension = flag.String("extension", "*", "Extension of files to import. This option can be used in order to only import RAWs or JPEGs.")
	_workers   = flag.Int("workers", runtime.NumCPU(), "The number of workers to use.")
	_filename  = flag.String(
		"filename",
		"{{.Date}}/{{.Date}}__{{.Time}}__{{.Name}}.{{.Extension}}",
		`
Template for the destination filename. The following properties are being passed to the template:

File Information:
- Name: 	 Original filename of the file excluding the extension
- Extension: Extension of the file without a leading dot.

Date information:
- Date: 	 Full date the image was taken (In the form YYYY-MM-DD, e.g. 2019-01-01)
- Time: 	 Full time the image was taken (In the form HH-MM-SS, e.g. 17-54-32)

More fine grained date information:
- Year:	 	 4-digit year the image was taken (e.g. 2019)
- Month:	 2-digit zero-padded month the image was taken (e.g. 01)
- Day: 		 2-digit zero-padded day the image was taken (e.g. 01)
- Hour: 	 2-digit zero-padded hour the image was taken (e.g. 17)
- Minute: 	 2-digit zero-padded minute the image was taken (e.g. 54)
- Second: 	 2-digit zero-padded second the image was taken (e.g. 32)

For more information about templating syntax, see the documentation at
https://golang.org/pkg/text/template
		`,
	)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, " Image Importer \n")
		fmt.Fprintf(os.Stderr, "===========\n")
		fmt.Fprintf(os.Stderr, "Simple tool to import images from an SD-Card into a defined\nfolder structure based on EXIF Information in the files.\n\n")

		fmt.Fprintf(os.Stderr, "image-importer will create an index file to track which files have already been imported.\n")
		fmt.Fprintf(os.Stderr, "This way only new images are being processed.\n\n")

		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "$ image-importer {OPTIONAL ARGUMENTS} --from /media/j/sd-card --to $HOME/Pictures/RAW \n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n\nExamples:\n\n")
		fmt.Fprintf(os.Stderr, `
$ image-importer --from /media/j/sd-card/DCIM/100OLYMP --to $HOME/Pictures/RAW
Importing images using 8 workers
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-39-49__P1010602.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-39-49__P1010602.JPG
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-20__P1010608.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-20__P1010608.JPG
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-07__P1010606.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-07__P1010606.JPG

2019/01/03 11:01:25 Imported 6 images`)
		fmt.Fprintf(os.Stderr, "\n\n")
		fmt.Fprintf(os.Stderr, `
$ image-importer --from /media/j/sd-card/DCIM/100OLYMP --to $HOME/Pictures/RAW --extension .ORF --filename "{{.Year}}/{{.Month}}/{{.Day}}/{{.OriginalName}}.{{.Extension}}"
Importing images using 8 workers
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010602.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010608.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010606.ORF

2019/01/03 11:01:25 Imported 3 images`)
	}
}

func parseSettings() Settings {
	flag.Parse()
	if *_from == "" || *_to == "" {
		log.Fatal("Options --from and --to are required.")
	}

	indexFile := fmt.Sprintf("%s/%s", *_from, *_index)

	lib.ExitIfError(lib.EnsureFileExists(indexFile))
	err, blacklist := lib.ReadLines(indexFile)
	lib.ExitIfError(err)

	files, err := filepath.Glob(fmt.Sprintf("%s/*%s", *_from, *_extension))
	lib.ExitIfError(err)

	return Settings{
		To:               *_to,
		From:             *_from,
		BlacklistPath:    indexFile,
		Blacklist:        blacklist,
		Workers:          *_workers,
		FilesRaw:         files,
		Files:            applyBlacklist(files, blacklist),
		FilenameTemplate: *template.Must(template.New("filenameTemplate").Parse(*_filename)),
	}
}

func applyBlacklist(files []string, blacklist []string) []string {

	filtered := make([]string, 0)
	for _, file := range files {

		isIgnored := false
		for _, ignore := range blacklist {
			if strings.Contains(file, ignore) {
				isIgnored = true
			}
		}

		if !isIgnored {
			filtered = append(filtered, file)
		}
	}
	return filtered
}
