package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("invalid command %v\n", os.Args[1])
	}
}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing %v\n", password)
		return
	}
	// convert the returned byte slice of numbers into a string
	hash := string(hashedBytes)
	fmt.Println(hash)
}

func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Invalid password: %v\n", password)
	}
	fmt.Println("Password is correct!")
}
