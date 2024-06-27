// Package internal contains internal logic of the application.
package pkg

import (
	"log"
	"os"
)

// RemoveFile removes the file with the specified path
func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Print(err)
	}
}
