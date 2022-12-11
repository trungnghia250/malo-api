package mail

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

func Send(personalization []*mail.Personalization, msgOpt ...OptMessage) (string, error) {
	msg := NewMessage()
	opts := []OptMessage{WithMessageEmailFrom("MALO CRM", "cskh.malo@gmail.com")}
	opts = append(opts, msgOpt...)
	msg.applyMessage(opts...)
	batchID := GenerateBatchID()
	msg.mailMessage.BatchID = batchID

	msg.mailMessage.AddPersonalizations(personalization...)
	request := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(msg.mailMessage)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return batchID, err
}

func GenerateBatchID() string {
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest("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", "/v3/mail/batch", host)
	request.Method = "POST"
	var result result
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(response.Body), &result)
	return result.BatchID
}

type result struct {
	BatchID string `json:"batch_id"`
}

func CancelScheduled(batchID string) error {
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest("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", "/v3/user/scheduled_sends", host)
	request.Method = "POST"
	request.Body = []byte(fmt.Sprintf(`{
  		"batch_id": "%s",
  		"status": "cancel"
	}`, batchID))
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	}
	return nil
}
