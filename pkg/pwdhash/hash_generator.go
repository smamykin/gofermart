package pwdhash

import "golang.org/x/crypto/bcrypt"

type HashGenerator struct {
}

func (h *HashGenerator) Generate(stringToHash string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)

	return string(b), err
}

func (h *HashGenerator) IsEqual(hashedPassword string, plainTxtPwd string) (isValid bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTxtPwd))
	if err == nil {
		return true, nil
	}

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}
