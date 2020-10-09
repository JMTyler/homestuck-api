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
// TODO: Should switch to using a struct like below, or StoryArc.Scrub(), instead of passing params directly.
func Ping(event string, story string, arc string, domain string, endpoint string, page int) {
	if fcmClient == nil {
		Init()
	}

	// import "homestuck/fcm"
	// fcm.SendMessage(map[string]string{
	// 	"endpoint": endpoint,
	//	"page":     fmt.Sprintf("%v", page),
	// })

	// payload := &struct {
	// 	Story    string
	// 	Arc      string
	// 	Endpoint string
	// 	Page     int
	// }{
	// 	Story:    story,
	// 	Arc:      arc,
	// 	Endpoint: endpoint,
	// 	Page:     page,
	// }
	// fmt.Println("Constructing payload:", payload)

	// v1
	if domain != "homestuck.com" {
		return
	}

	// TODO: See documentation on defining a message payload.
	message := &messaging.Message{
		Topic: "v1-" + topic,
		Data: map[string]string{
			"event":    event,
			"endpoint": endpoint,
			"title":    story,
			"subtitle": arc,
			"pages":    fmt.Sprintf("%v", page),
		},
	}

	_, err := fcmClient.Send(getContext(), message)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully sent FCM message.")
}

func Subscribe(version string, token string) error {
	if fcmClient == nil {
		Init()
	}

	_, err := fcmClient.SubscribeToTopic(getContext(), []string{token}, version+"-"+topic)
	if err != nil {
		return err
	}
	// fmt.Println(response.SuccessCount, "tokens were subscribed to", topic)

	return nil
}

func Unsubscribe(version string, token string) error {
	if fcmClient == nil {
		Init()
	}

	_, err := fcmClient.UnsubscribeFromTopic(getContext(), []string{token}, version+"-"+topic)
	if err != nil {
		return err
	}
	// fmt.Println(response.SuccessCount, "tokens were unsubscribed from", topic)

	return nil
}
