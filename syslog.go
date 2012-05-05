package main
/*
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <syslog.h>

// Wrapper for the va_list version of syslog.
void _syslog(int priority, const char *msg) {
	va_list l;
	vsyslog(priority, msg, l);
}
*/
import "C"

import (
	"log/syslog"
	"unsafe"
)

const (
	LOG_NDELAY = 0x08
	LOG_LOCAL1 = 17 << 3
)

// Identity of the currently open syslog.  Zero value is not opened yet
var syslogIdent string

// Single threaded
func Syslog(ident string, priority syslog.Priority, message string) {
	// (re)open the log
	if ident != syslogIdent {
		if syslogIdent != "" {
			C.closelog()
		}

		C.openlog(C.CString(ident), LOG_NDELAY, LOG_LOCAL1)
		syslogIdent = ident
	}

	s := C.CString(message)
	C._syslog(C.int(priority), s)
	C.free(unsafe.Pointer(s))
}
