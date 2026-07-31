package auth

import "golang.org/x/crypto/bcrypt"

// ComparePassword compara un hash bcrypt
// contra una contraseña recibida.
//
func ComparePassword(hash string, password string) error {

	return bcrypt.CompareHashAndPassword(

		[]byte(hash),

		[]byte(password),

	)

}
