package auth

import "golang.org/x/crypto/bcrypt"

// Hash Plain Text password from request, before storing in DB.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare plain text password in request and hashed password from record...
// return true if both are same, false otherwise.
func ComparePasswords(hashed string, plainText []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plainText)
	return err == nil
}
