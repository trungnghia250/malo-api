package mail

import "github.com/sendgrid/sendgrid-go/helpers/mail"

type OptMessage func(*Message)

func WithMessageEmailFrom(name string, address string) OptMessage {
	return func(message *Message) {
		message.mailMessage.From = &mail.Email{
			Name:    name,
			Address: address,
		}
	}
}

func WithMessageHTML(messageContent string) OptMessage {
	return func(message *Message) {
		content := &mail.Content{
			Type:  "text/html",
			Value: messageContent,
		}
		message.mailMessage.Content = append(message.mailMessage.Content, content)
	}
}

func NewMessage() *Message {
	return &Message{
		mailMessage: mail.NewV3Mail(),
	}
}

type Message struct {
	mailMessage *mail.SGMailV3
}

func (m *Message) applyMessage(opts ...OptMessage) {
	for _, opt := range opts {
		opt(m)
	}
}
