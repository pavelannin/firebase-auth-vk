package firebase

import "fmt"

// CreateCustomToken -- Создает Firebase auth токен из пользовательского идентификатора
func (env *Environment) CreateCustomToken(userUID string) (string, error) {
	token, err := env.Auth.CustomToken(userUID)
	if err != nil {
		return "", fmt.Errorf("[Firebase] error minting custom token: %v", err)
	}
	return token, nil
}
