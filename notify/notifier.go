// Package Notify is used for sending email/sms notifications
package notify

import (
	"fmt"
	"io"
	"net/smtp"
	"os"
)

type Notifier interface {
	io.Writer
}

type MockNotifier struct{}

// Write sends bytes to stdout
func (n MockNotifier) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

type EmailNotifier struct{}

func (n EmailNotifier) Write(p []byte) (int, error) {
	// Set up authentication information.
	f := os.Getenv("FROM_EMAIL")
	t := os.Getenv("TO_EMAIL")
	ps := os.Getenv("EMAIL_PASS")
	auth := smtp.PlainAuth("", f, ps, "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{t}
	msg := "To: " + string(t) + "\r\nSubject: " + "Prayzo Monitoring Alert\r\n\r\n" + string(p)
	err := smtp.SendMail("smtp.gmail.com:587", auth, f, to, []byte(msg))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
