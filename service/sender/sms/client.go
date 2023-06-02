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
	accountSid := "AC74acd33e928f3b9cfdfba86132a51d2e"
	authToken := "eefe57305498e33688a3ce0fd253377e"
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &api.CreateMessageParams{}
	params.SetFrom("+13612822731")
	params.SetBody(request.Content)
	if request.SendAt > 0 {
		params.SetSendAt(time.Unix(int64(request.SendAt), 0))
		params.SetMessagingServiceSid("MG4fe28942242424956b087cd4794a3333")
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
	accountSid := "AC74acd33e928f3b9cfdfba86132a51d2e"
	authToken := "eefe57305498e33688a3ce0fd253377e"
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
