package firebase

import (
	"context"
	"fmt"
	"sync"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	app  *firebase.App
	once sync.Once
)

func InitializeFirebase(credentialsFile string) (*firebase.App, error) {
	var err error
	once.Do(func() {
		opt := option.WithCredentialsFile(credentialsFile)
		app, err = firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			err = fmt.Errorf("error initializing app: %v", err)
		}
	})
	return app, err
}

func GetFirebaseApp() *firebase.App {
	return app
}
