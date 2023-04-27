package main

import (
	"fmt"
	"log"
	"os"
)

func directly(message string) {
	// write a message to stderr
	fmt.Fprint(os.Stderr, message)
}

func createLogger() *log.Logger {
	// https://www.honeybadger.io/blog/golang-logging/
	// You can use log flag constants to enrich a log message
	// by providing additional context information,
	// such as the file, line number, date, and time.
	// For example, passing the message "Something went wrong" through a logger with a flag combination shown below:
	// log.Ldate|log.Ltime|log.Lshortfile (bit mask to enable fields in output)
	return log.New(os.Stderr, "oops! - ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	directly(fmt.Sprintf("the error message is %q\n", "d'oh!"))
	l := createLogger()
	l.Println("Your music's bad and you should feel bad!")
}
