#include <string.h>
#include <stdlib.h>

typedef struct {
    char* path;
    Date date;

    ExifData *exif_data;
} File;

File create_file(char* path) {
    File file;
    file.path = path;
    file.exif_data = exif_data_new_from_file(path);
    
    ExifTag possible_date_tags[] = { EXIF_TAG_DATE_TIME_ORIGINAL, EXIF_TAG_DATE_TIME, EXIF_TAG_DATE_TIME_DIGITIZED };
    ExifEntry* entry;
    for (int i = 0; i < sizeof(possible_date_tags); i++) {
        entry = exif_data_get_entry(file.exif_data, possible_date_tags[i]);
        if (entry != NULL) {
            file.date = parse_date(entry->data);
            break;
        }
    }

    return file;
}

void print_file(File file) {
    printf("path=%s \n", file.path);
    printf("date.year=%i \n", file.date.year);
    printf("date.month=%i \n", file.date.month);
    printf("date.day=%i \n", file.date.day);
    printf("date.minute=%i \n", file.date.minute);
    printf("date.hour=%i \n", file.date.hour);
    printf("date.second=%i \n", file.date.second);
}

// Copied from https://stackoverflow.com/questions/779875/what-is-the-function-to-replace-string-in-c/779960#779960
// TODO understand
char *str_replace(char *orig, char *rep, char *with) {
    char *result; // the return string
    char *ins;    // the next insert point
    char *tmp;    // varies
    int len_rep;  // length of rep (the string to remove)
    int len_with; // length of with (the string to replace rep with)
    int len_front; // distance between rep and end of last rep
    int count;    // number of replacements

    // sanity checks and initialization
    if (!orig || !rep)
        return NULL;
    len_rep = strlen(rep);
    if (len_rep == 0)
        return NULL; // empty rep causes infinite loop during count
    if (!with)
        with = "";
    len_with = strlen(with);

    // count the number of replacements needed
    ins = orig;
    for (count = 0; tmp = strstr(ins, rep); ++count) {
        ins = tmp + len_rep;
    }

    tmp = result = malloc(strlen(orig) + (len_with - len_rep) * count + 1);

    if (!result)
        return NULL;

    // first time through the loop, all the variable are set correctly
    // from here on,
    //    tmp points to the end of the result string
    //    ins points to the next occurrence of rep in orig
    //    orig points to the remainder of orig after "end of rep"
    while (count--) {
        ins = strstr(orig, rep);
        len_front = ins - orig;
        tmp = strncpy(tmp, orig, len_front) + len_front;
        tmp = strcpy(tmp, with) + len_with;
        orig += len_front + len_rep; // move to next "end of rep"
    }
    strcpy(tmp, orig);
    return result;
}


char* format_new_path(File file, char* template) {
    char year[5];
    sprintf(year, "%04d", file.date.year);
    char month[3];
    sprintf(month, "%02d", file.date.month);
    char day[3];
    sprintf(day, "%02d", file.date.day);
    char hour[3];
    sprintf(hour, "%02d", file.date.hour);
    char minute[3];
    sprintf(minute, "%02d", file.date.minute);
    char second[3];
    sprintf(second, "%02d", file.date.second);


    char* new_name = str_replace(template, "{year}", year);
    new_name = str_replace(new_name, "{month}", month);
    new_name = str_replace(new_name, "{day}", day);
    new_name = str_replace(new_name, "{hour}", hour);
    new_name = str_replace(new_name, "{minute}", minute);
    new_name = str_replace(new_name, "{second}", second);

    return new_name;
}