package mail

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

func Send(personalization []*mail.Personalization, msgOpt ...OptMessage) error {
	msg := NewMessage()
	opts := []OptMessage{WithMessageEmailFrom("MALO CRM", "cskh.malo@gmail.com")}
	opts = append(opts, msgOpt...)
	msg.applyMessage(opts...)

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
	return err
}
