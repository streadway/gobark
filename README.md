# Bark

Pipes and formats lines from stdin to /dev/log and optionally tees the raw output back out to stdout.

Intended to be used with services under daemontools or runit as a log tee with some payload introspection.

This should be able to sustain 100k lines/s, at about 20MB/s.

GC is a concern and the observed latency is about 1ms for a cycle, with a GC cycle happening around every second.  For latency sensitive programs that log unbuffered to stdout, please test carefully.

# Usage

    Usage of ./gobark:
      -delim="\n": the byte sequence that separates events
      -ignore-delim=false: exclude the delimiter in each event
      -line-buffers=1000: max number of events to buffer
      -line-size=4096: expect most events to be shorter than this size
      -name="bark": identity/program name that is prefixed to each event
      -tee=false: immediately write and block on all reads to stdout
      -xpid=false: checks each event for the [x-pid=""] header

Example

    ./myserver | ./gobark -name myserver -tee > /var/log/myserver

If your program emits lines with the syslog extension `[ x-pid="program" ...`, bark will find it and use `program` as the program name instead of the one passed on the command line.

# Testing

Build the timings package with `make` then collect and bucket the latencies:

    cd timing
    make
    cd ..
    out=samples.out
    timing/timing 2>$out | gobark
    cat $out | awk '{print int(log($0)/log(10))}' |\
      sort -n | uniq -c | awk '{print $2, $1}' | sort -nr

# Contributing

Make your changes on your fork in a non-master branch and submit a pull request with the use case/bug you've contributed to helping with.

Please run `go fmt` before submitting your patch.

Check the issues tab in github for any ideas if you'd like to contribute.

# Credits

Omid Aladini - for the initial implementation

# License

Copyright (c) 2012, Sean Treadway, Omid Aladini

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

Redistributions in binary form must reproduce the above copyright notice, this
list of conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
