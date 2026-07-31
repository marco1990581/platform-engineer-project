package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {

	fmt.Println("===================================")
	fmt.Println(" Platform API - Create User")
	fmt.Println("===================================")
	fmt.Println()

	// Pedimos el nombre de usuario.
	fmt.Print("Username: ")

	var username string
	fmt.Scanln(&username)

	// Pedimos la contraseña sin mostrarla.
	fmt.Print("Password: ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	// Confirmación.
	fmt.Print("Confirm password: ")

	confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	// Verificamos que ambas coincidan.
	if string(password) != string(confirmPassword) {

		log.Fatal("Passwords do not match")

	}

	// Generamos el hash bcrypt.
	hash, err := bcrypt.GenerateFromPassword(
		password,
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("User created successfully!")
	fmt.Println()

	fmt.Println("Username:")
	fmt.Println(username)

	fmt.Println()

	fmt.Println("Password Hash:")
	fmt.Println(string(hash))

}
