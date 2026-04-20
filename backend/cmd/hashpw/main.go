package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"alice123", "bob123", "admin123"}

	for _, p := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s -> %s\n", p, string(hash))
	}
}