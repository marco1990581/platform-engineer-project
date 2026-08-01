package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marcosalbano/platform-api/internal/auth"
	"github.com/marcosalbano/platform-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {

	fmt.Println("===================================")
	fmt.Println(" Platform API - Create User")
	fmt.Println("===================================")
	fmt.Println()

	var username string

	fmt.Print("Username: ")
	fmt.Scanln(&username)

	fmt.Print("Password: ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	fmt.Print("Confirm password: ")

	confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	if string(password) != string(confirmPassword) {
		log.Fatal("passwords do not match")
	}

	hash, err := bcrypt.GenerateFromPassword(
		password,
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Fatal(err)
	}

	repository := auth.NewFileRepository("configs/users.json")

	user := models.User{

		Username: username,

		PasswordHash: string(hash),

		Role: "admin",
	}

	err = repository.AddUser(user)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("✔ User created successfully")
}
