package beanimplement

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"google.golang.org/api/option"
	"os"
)

type NotificationClient struct {
	client *messaging.Client
}

func getDecodedFireBaseKey() ([]byte, error) {

	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}

func NewNotificationClient() bean.NotificationClient {
	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return nil
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	// Initialize firebase app
	app, err := firebase.NewApp(context.Background(), nil, opts...)

	if err != nil {
		fmt.Println("Error in initializing firebase app: %s", err)
		return nil
	}

	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return nil
	}

	return &NotificationClient{
		client: fcmClient,
	}
}

func (c *NotificationClient) Send() {
	_, err := c.client.Send(context.Background(), &messaging.Message{

		Notification: &messaging.Notification{
			Title: "Test Backend!!",
			Body:  "test backend",
		},
		Token: "dPxTJmpPR_UwN6qizbMSFC:APA91bF1M3fTPC4Gex2N5vZ_Yx2Gr44OQJzQBnyt2MXCDHLua_Z5rBl2eh-EEgzmsictZAATJjtWGn0nCM9zyGk0U5saXVDI_H91I9wM7HHso3G81SVsgZU", // it's a single device token
	})

	if err != nil {
		fmt.Println("Error sending message: %s", err)
	}
}

func (c *NotificationClient) SendTest(req model.NotificationRequest) {
	_, err := c.client.Send(context.Background(), &messaging.Message{

		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Token: req.Token, // it's a single device token
	})

	if err != nil {
		fmt.Println("Error sending message: %s", err)
	}
}
