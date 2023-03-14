package pwdhash

import "golang.org/x/crypto/bcrypt"

type HashGenerator struct {
}

func (h *HashGenerator) Generate(stringToHash string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)

	return string(b), err
}
