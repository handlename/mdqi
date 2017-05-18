package mdqi

import (
	"log"
	"os"
)

var logger = func() *log.Logger {
	return log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}()
