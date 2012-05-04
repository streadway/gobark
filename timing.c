// test harness that will measure the high resolution write times on OSX across
// 100k log lines at full saturation
//
// g++ -framework CoreServices -lc -o timings timing.c
//
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
