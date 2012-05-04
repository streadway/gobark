/***
 test harness that will measure the high resolution write times on OSX across
 100k log lines at full saturation

 g++ -framework CoreServices -lc -o timings timing.c

 Run the timing benchmark and write out the raw sample times

 ./timings 2>samples | ./gobark

 Then chew on them.  Here's a log(10) histogram bashness to identify how variant the outliers are:

 f=samples; cat $f |\
  awk '{print int(log($0)/log(10))}' |\
  sort -n |\
  uniq -c |\
  awk '{print "10^" $2, $1}' |\
  sort -nr > $f.log10

***/
#include <stdio.h>
#include <math.h>
#include <CoreServices/CoreServices.h>

const char *line = "12345678901234567890123456789012345678901234567890\n";
const int iterations = 100000;

void bark(const char *line, const int len) {
    int n = fwrite(line, 1, len, stdout);
    if (n < len) {
        printf("short write: %d %d\n", n, len);
        exit(1);
    }
}

int main(int argc, char **argv) {
    int timings[iterations];
    int i;

    int len = strlen(line);

    //Set output buffer to zero
    setvbuf (stdin, NULL, _IONBF, 0);
    setvbuf (stdout, NULL, _IONBF, 0);

    for (i = 0; i < iterations; i++) {
        AbsoluteTime start = UpTime();

        bark(line, len);

        AbsoluteTime end = UpTime();

        Nanoseconds diffNS = AbsoluteDeltaToNanoseconds(end,start);
        timings[i] = UnsignedWideToUInt64(diffNS);
    }

    for (i =0; i < iterations; i++) {
        fprintf(stderr, "%d\n", timings[i]);
    }

    return 0;
}
