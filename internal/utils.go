package internal

import (
	"log"
	"os"
)

// RemoveFile removes the file with the specified path
func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Print(err)
	}
}
