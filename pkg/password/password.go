package password

import "golang.org/x/crypto/bcrypt"

// Hash generates a bcrypt hash from plaintext password.
func Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Compare checks whether a plaintext password matches a bcrypt hash.
func Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
