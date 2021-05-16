package main

import (
	"crypto/rand"
	"fmt"
	"log"
)

func generateGUID() (s string) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x",
		b[0:4], b[4:6], b[6:])
	return uuid
}

func generatePassword(size int) (s string) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x",
		b[0:size])
	return uuid
}

func boolToInt(boolToCheck bool) (i int) {
	if boolToCheck {
		return 1
	} else {
		return 0
	}
}
