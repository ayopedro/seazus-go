package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash), err
}

func comparePasswords(pw, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err
}
