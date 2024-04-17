package internal

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

// RemoveFile removes the file with the specified path
func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal()
	}
}

// CheckForErrorAndSendStatusWithLog if an error exists sends the specified status and logs the error
func CheckForErrorAndSendStatusWithLog(c *fiber.Ctx, err error, status int) {
	if err != nil {
		_ = c.SendStatus(status)
		log.Fatal(err)
	}
}
