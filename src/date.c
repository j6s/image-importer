
typedef struct {
    int year;
    int month;
    int day;
    int hour;
    int minute;
    int second;
} Date;

Date parse_date(char* date_string) {
    // Incoming dates have the format YYYY:MM:DD HH:ii:ss
    Date date;

    sscanf(
        date_string, 
        "%d:%d:%d %d:%d:%d", 
        &date.year, 
        &date.month, 
        &date.day, 
        &date.hour, 
        &date.minute, 
        &date.second
    );

    return date;
}