package email

import "fmt"

type Email struct {
	MessageID string   `json:"messageID"`
	Date      string   `json:"date"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
}

func (e Email) ToString() string {
	return fmt.Sprintf("Message ID: %v \n Date: %v \n From: %v \n To: %v \n Subject: %v \n Content: %v \n ", e.MessageID, e.Date, e.From, e.To, e.Subject, e.Content)
}
