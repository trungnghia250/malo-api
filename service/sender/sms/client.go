package sms

import (
	"fmt"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"time"
)

type Request struct {
	Receivers []string
	Content   string
	SendAt    int32
}

func Send(request Request) ([]string, error) {
	accountSid := "AC3875387c245514321b16003f0f109794"
	authToken := "736de2bfadf1fe37ad6f4cfc6ea201e3"
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &api.CreateMessageParams{}
	params.SetFrom("+15635002841")
	params.SetBody(request.Content)
	if request.SendAt > 0 {
		params.SetSendAt(time.Unix(int64(request.SendAt), 0))
		params.SetMessagingServiceSid("MG28c99338dbe0159c1a609d0e9f2865db")
		params.SetScheduleType("fixed")
	}
	var messageSid []string
	for _, receiver := range request.Receivers {
		params.SetTo(receiver)
		resp, err := client.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		}
		messageSid = append(messageSid, *resp.Sid)
	}

	return messageSid, nil
}

func Cancel(messageSIDs []string) error {
	accountSid := "AC3875387c245514321b16003f0f109794"
	authToken := "736de2bfadf1fe37ad6f4cfc6ea201e3"
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &api.UpdateMessageParams{}
	params.SetStatus("canceled")
	for _, messageSID := range messageSIDs {
		_, err := client.Api.UpdateMessage(messageSID, params)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
