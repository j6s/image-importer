#include <stdio.h>
#include <libexif/exif-data.h>
#include "date.c"
#include "file.c"

int main(int argc, char* argv[]) {

    // TODO read command line flags

    // TODO Read list of files from directory

    // TODO actually do the copying

    // TODO read and write blacklist file
    
    for (int i = 1; i < argc; i++) {
        printf("[%d]\n", i);
        File file = create_file(argv[i]);
        printf("%s -> %s \n", file.path, format_new_path(file, "{year}-{month}-{day}/{hour}-{minute}-{second}"));
    }


    return 0;
}