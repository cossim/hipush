package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("/Users/macos-15/Downloads/cossim-5a21a-firebase-adminsdk-atk43-5a6cdebc4e.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	//client, err := app.Auth(ctx)
	//if err != nil {
	//	log.Fatalf("error getting Auth client: %v\n", err)
	//}
	//
	//token, err := client.CustomToken(ctx, "some-uid")
	//if err != nil {
	//	log.Fatalf("error minting custom token: %v\n", err)
	//}

	//log.Printf("Got custom token: %v\n", token)

	send(ctx, app)
}

func send(ctx context.Context, app *firebase.App) {
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "dbogAnZeQ5iM85zVMncB_8:APA91bHP956b7LUtbSJTOHMvlukLuU9uliFcHvsTTp5i50azaMjM2VfkAZ1zC46FM9T9RDjV9JH-F1qlrGD97mobsNPxu8RPwoD1sRfi1r2TgB74lnw5DZGdSrKsZebKnUyMF0BIJ6hv"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		//Data: map[string]string{
		//	"score": "850",
		//	"time":  "2:45",
		//},
		Notification: &messaging.Notification{
			Title:    "测试标题",
			Body:     "测试内容",
			ImageURL: "",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}
