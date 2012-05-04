// Tees your events from stdin into syslog
package main

import (
	"bytes"
	"os"
	"flag"
	"regexp"
	"log/syslog"
)

// Accumulators for the different log variants we could see
type loggers map[string]*syslog.Writer
type priorities map[syslog.Priority]loggers

var (
	pool chan *bytes.Buffer
	lines chan *bytes.Buffer

	logs = make(priorities)
	pid = regexp.MustCompile(`\[[^]]*x-pid="([^"]+)"`)
)

var (
	name = flag.String("name", "bark", "identity/program name that is prefixed to each event")
	xpid = flag.Bool("xpid", false, "checks each event for the [x-pid=\"\"] header")
	tee = flag.Bool("tee", false, "immediately write and block on all reads to stdout")
	delim = flag.String("delim", "\n", "the byte sequence that separates events")
	ignoreDelim = flag.Bool("ignore-delim", false, "exclude the delimiter in each event")
	lineBuffers = flag.Int("line-buffers", 1000, "max number of events to buffer")
	lineSize = flag.Int("line-size", 4096, "expect most events to be shorter than this size")
)

func init() {
	flag.Parse()

	// reusable byte buffers to reduce churn
	pool = make(chan *bytes.Buffer, *lineBuffers)
	lines = make(chan *bytes.Buffer, *lineSize)

	for i := 0; i < *lineBuffers; i++ {
		pool <- new(bytes.Buffer)
	}
}

func free(b *bytes.Buffer) {
	b.Reset()
	pool <- b
}

func reader() {
	scratch := make([]byte, *lineSize)

	for {
		n, err := os.Stdin.Read(scratch[:])
		if err != nil {
			close(lines)
			return
		}

		part := scratch[:n]

		// tee time, two lumps please
		if *tee {
			if _, err = os.Stdout.Write(part); err != nil {
				close(lines)
				return
			}
		}

		for len(part) > 0 {
			line := <-pool
			i := bytes.Index(part, []byte(*delim))

			if i < 0 {
				// partial write, will get aggregated in main()
				line.Write(part)
				part = part[0:0]
			} else {
				// full write, strip the delim in main()
				line.Write(part[:i+len(*delim)])
				part = part[i+len(*delim):]
			}

			lines <- line
		}
	}
}

// Looks for x-pid for the name or the configured name
func prefix(line []byte) string {
	if *xpid {
		if matches := pid.FindSubmatch(line); matches != nil {
			return string(matches[1])
		}
	}

	return *name
}

func priority(line []byte) (prio syslog.Priority, clean []byte) {
	var start string

	prio, clean = syslog.LOG_INFO, line

	if i := bytes.IndexByte(line, ' '); i >= 0 {
		start, clean = string(line[:i]), line[i+1:]

		switch start {
		case "EMERG", "EMERGENCY":
			prio = syslog.LOG_EMERG
  	case "ALERT":
			prio = syslog.LOG_ALERT
  	case "CRIT", "CRITICAL":
			prio = syslog.LOG_CRIT
  	case "ERROR", "ERR":
			prio = syslog.LOG_ERR
  	case "WARN", "WARNING":
			prio = syslog.LOG_WARNING
  	case "NOTICE":
    	prio = syslog.LOG_NOTICE
		case "INFO":
    	prio = syslog.LOG_INFO
		case "DEBUG":
    	prio = syslog.LOG_DEBUG
		}
	}

	return
}

// creates and registers a new logger at the prefix from our
// name or x-pid
func logger(line []byte) (log *syslog.Writer, clean []byte, err error) {
	prio, clean := priority(line)
	name := prefix(clean)

	cache, ok := logs[prio]
	if !ok {
		cache = make(loggers)
		logs[prio] = cache
	}

	log, ok = cache[name]
	if !ok {
		if log, err = syslog.New(prio, name); err != nil {
			return
		}
		cache[name] = log
	}

	return
}

// process and bark out a line, stripping the delimiter if needs be
func bark(line []byte) error {
	if *ignoreDelim {
		line = line[:len(line)-len(*delim)]
	}

	log, line, err := logger(line)
	if err != nil {
		return err
	}
	
	log.Write(line)

	return err
}

func main() {
	var line bytes.Buffer

	go reader()

	for {
		line.Reset()

		// Combine the next possible parts to a complete line
		for !bytes.HasSuffix(line.Bytes(), []byte(*delim)) {
			next, ok := <-lines

			if !ok {
				// reader closed, graceful shutdown
				os.Exit(0)
			}

			line.Write(next.Bytes())
			free(next)
		}

		if err := bark(line.Bytes()); err != nil {
			os.Exit(1)
		}
	}
}

// Copyright (c) 2012, Sean Treadway, Omid Aladini
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
// Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
// 
// Redistributions in binary form must reproduce the above copyright notice, this
// list of conditions and the following disclaimer in the documentation and/or
// other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
