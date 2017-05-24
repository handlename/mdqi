package mdqi

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var logger = func() *log.Logger {
	return log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}()

var debug = func() *log.Logger {
	var out io.Writer

	if os.Getenv("DEBUG") != "" {
		out = os.Stderr
	} else {
		out = ioutil.Discard
	}

	l := log.New(out, "[debug] ", log.Ldate|log.Ltime|log.Lshortfile)

	return l
}()