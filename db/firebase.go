package db

import (
	"context"
	"log"

	cloud "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	Storage *cloud.Client
)

func InitFirebaseApp() error {
	sa := option.WithCredentialsFile("ezcom-eaa21-6af33392e490.json")

	var err error

	Storage, err = cloud.NewClient(context.Background(), sa)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
