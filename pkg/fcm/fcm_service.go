package fcm

import (
	"context"
	"fmt"
	"log"
	"os"

	"firebase.google.com/go/v4/messaging"
	"github.com/appleboy/go-fcm"
)

var fcmClient *fcm.Client

func InitializeFCM() {
	var err error

	fcmClient, err = fcm.NewClient(
		context.TODO(),
		fcm.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_PATH")),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func SendFCMAlert(clientId string, payload string) {
	token := os.Getenv("FCM_ANDROID_TARGET_ID_TOKEN")
	resp, err := fcmClient.Send(
		context.TODO(),
		&messaging.Message{
			Token: token,
			Data: map[string]string{
				// "message": fmt.Sprintf("{\"datetime\": \"%d\", \"content\": \"Danger in front\",\"alert_type\": 1 }", time.Now().Unix()),
				"message": payload,
			},
			Android: &messaging.AndroidConfig{
				Priority: "high",
			},
			APNS: &messaging.APNSConfig{
				Headers: map[string]string{
					"apns-priority": "10",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success count:", resp.SuccessCount)
	fmt.Println("failure count:", resp.FailureCount)
	fmt.Println("message id:", resp.Responses[0].MessageID)
	fmt.Println("error msg:", resp.Responses[0].Error)
}
