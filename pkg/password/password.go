package password

import "golang.org/x/crypto/bcrypt"

// Hash generates a bcrypt hash from plaintext password with default cost.
func Hash(plain string) (string, error) {
	return HashWithCost(plain, bcrypt.DefaultCost)
}

// HashWithCost generates a bcrypt hash from plaintext password with a specified cost factor.
func HashWithCost(plain string, cost int) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Compare checks whether a plaintext password matches a bcrypt hash.
func Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
