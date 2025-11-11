package handler

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587 // default for Gmail
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("❌ Email send error:", err)
		return err
	}

	fmt.Printf("📨 Email sent successfully to %s\n", to)
	return nil
}
