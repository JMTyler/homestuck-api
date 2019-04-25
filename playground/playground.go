package main

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"golang.org/x/net/context"
)

func main() {
	fmt.Println("Hello, playground")

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "eT2xBsqTLVc:APA91bFlGhfSTMw4hWj5apVCM8YuNJVxl7Osg9GhVbP_lSqNYqPYMKn2YteRy9eV4tHfbtHriER39Nq1QHAKPYGO3_5yjjr3EclUkK6rX3t78lveCUxQVSQubfyGBA5Y5fXILqLjcNuV"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		// Token, Topic, or Condition
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		panic(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}
