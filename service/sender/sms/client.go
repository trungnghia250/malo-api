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
	SendAt    time.Time
}

func Send(request Request) error {
	accountSid := "AC3875387c245514321b16003f0f109794"
	authToken := "736de2bfadf1fe37ad6f4cfc6ea201e3"
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &api.CreateMessageParams{}
	params.SetFrom("+15635002841")
	params.SetBody(request.Content)
	params.SetSendAt(request.SendAt)

	for _, receiver := range request.Receivers {
		params.SetTo(receiver)
		resp, err := client.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if resp.Sid != nil {
				fmt.Println(*resp.Sid)
			} else {
				fmt.Println(resp.Sid)
			}
		}
	}

	return nil
}
