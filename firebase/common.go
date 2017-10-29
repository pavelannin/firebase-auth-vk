package firebase

import (
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// Environment -- Стуктура содержит окружение взаимодействия с Firebase
type Environment struct {
	Auth *auth.Client
}

// New -- Создает экземпляр с окружением Firebase
func New(credentialsFile string) *Environment {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("[Firebase] error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("[Firebase] error initializing app: %v\n", err)
	}
	return &Environment{
		Auth: client,
	}
}
