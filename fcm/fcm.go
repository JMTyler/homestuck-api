package fcm

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"
)

const (
	PotatoEvent = "Potato"
	SyncEvent   = "SyncStory"
)

// TODO: Should also make this topic v1-specific if we end up prefixing the endpoints with /v1 as well.
const topic = "Stories"

var fcmClient *messaging.Client

func getContext() context.Context {
	return context.Background()
}

// HACK: Had to make this public. Why does Ping() not recognise this function if it's private?
func Init() {
	googleCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	app, err := firebase.NewApp(getContext(), nil, option.WithCredentialsJSON([]byte(googleCredentials)))
	if err != nil {
		panic(err)
	}

	fcmClient, err = app.Messaging(getContext())
	if err != nil {
		panic(err)
	}
}

// TODO: Publish
func Ping(event string, story string, arc string, endpoint string, page int) {
	if fcmClient == nil {
		Init()
	}

	// import "homestuck/fcm"
	// fcm.SendMessage(map[string]string{
	// 	"endpoint": endpoint,
	//	"page":     fmt.Sprintf("%v", page),
	// })

	payload := &struct {
		Story    string
		Arc      string
		Endpoint string
		Page     int
	}{
		Story:    story,
		Arc:      arc,
		Endpoint: endpoint,
		Page:     page,
	}
	fmt.Println("Constructing payload:", payload)

	// TODO: See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"event":    event,
			"story":    story,
			"arc":      arc,
			"endpoint": endpoint,
			"page":     fmt.Sprintf("%v", page),
		},
		// Token, Topic, or Condition
		Topic: topic,
	}

	_, err := fcmClient.Send(getContext(), message)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully sent FCM message.")
}

func Subscribe(registrationTokens []string) error {
	if fcmClient == nil {
		Init()
	}

	response, err := fcmClient.SubscribeToTopic(getContext(), registrationTokens, topic)
	if err != nil {
		return err
	}
	fmt.Println(response.SuccessCount, "tokens were subscribed to", topic)

	return nil
}

func Unsubscribe(registrationTokens []string) error {
	if fcmClient == nil {
		Init()
	}

	response, err := fcmClient.UnsubscribeFromTopic(getContext(), registrationTokens, topic)
	if err != nil {
		return err
	}
	fmt.Println(response.SuccessCount, "tokens were unsubscribed from", topic)

	return nil
}
