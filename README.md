# Image Importer

Simple tool to import images from an SD-Card into a defined folder structure based on EXIF Information in the files.
`image-importer` keeps track of images that have already been imported in order to not process the same image twice.
This way you can fill up your SD-Card without having to copy the same files every time - only new images will be processed.

This tool assumes a unixy-system (such as Linux, BSD or macOS) with [`exiftool`](https://sno.phy.queensu.ca/~phil/exiftool/) installed.

## Installation
### Step 1: Ensure `exiftool` is installed on your system
Debian / Ubuntu:
```bash
sudo apt install exiftool
```

macOS:
```bash
brew install exiftool
```

### Step 2: Install `image-importer`

- Option 1: Download a [compiled release binary](https://github.com/j6s/image-importer/releases) (No `go` needed)
- Option 2: Use `go get`
    ```bash
    go get github.com/j6s/image-importer/cmd/image-importer
    ```

## Usage

```
$ image-importer {OPTIONAL ARGUMENTS} --from /media/j/sd-card --to $HOME/Pictures/RAW 

Arguments:
  -extension string
    	Extension of files to import. This option can be used in order to only import RAWs or JPEGs. (default "*")
  -filename string
    	
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
    			 (default "{{.Date}}/{{.Date}}__{{.Time}}__{{.Name}}.{{.Extension}}")
  -from string
    	The source directory of the files to import. All files in this directory (matching the given extension) will be processed.
  -index string
    	Name of the index file in the source folder. (default ".imported")
  -to string
    	The destination directory of the files. The files will be copied to this folder.
  -workers int
    	The number of workers to use. (default 8)
```

## Examples

```
$ image-importer --from /media/j/sd-card/DCIM/100OLYMP --to $HOME/Pictures/RAW
Importing images using 8 workers
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-39-49__P1010602.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-39-49__P1010602.JPG
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-20__P1010608.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-20__P1010608.JPG
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.ORF -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-07__P1010606.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.JPG -> /home/j/Pictures/TEST/2019-01-01/2019-01-01__14-42-07__P1010606.JPG

2019/01/03 11:01:25 Imported 6 images
```

```
$ image-importer --from /media/j/sd-card/DCIM/100OLYMP --to $HOME/Pictures/RAW --extension .ORF --filename "{{.Year}}/{{.Month}}/{{.Day}}/{{.OriginalName}}.{{.Extension}}"
Importing images using 8 workers
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010602.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010602.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010608.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010608.ORF
2019/01/03 11:01:25 [---] /media/j/sd-card/DCIM/100OLYMP/P1010606.ORF -> /home/j/Pictures/TEST/2019/01/01/P1010606.ORF

2019/01/03 11:01:25 Imported 3 images
```

## TODO

- Use platform-agnostic copy method instead of shelling out to `cp`
- Add option to manually specify location of `exiftool` binary