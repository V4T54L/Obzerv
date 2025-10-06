package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	logModes = []string{"INFO", "WARNING", "ERROR", "DEBUG"}
	services = []string{"auth", "payment", "user", "inventory", "logging"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateLogEntry() string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMode := logModes[rand.Intn(len(logModes))]
	service := services[rand.Intn(len(services))]
	content := randomString(20)

	return fmt.Sprintf("%s, %s, %s, %s", timestamp, logMode, service, content)
}
