package mailer

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendVerificationEmail(toEmail string, token string, appURL string) error
}

type mailer struct {
	host string
	port int
	user string
	pass string
}

func NewMailer(host string, port int, user, pass string) Mailer {
	return &mailer{host: host, port: port, user: user, pass: pass}
}

func (m *mailer) SendVerificationEmail(toEmail string, token string, appURL string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.user)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Verify Your Email Address")

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", appURL, token)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
			<h2 style="color: #333;">Welcome to GOeShop!</h2>
			<p style="color: #555;">Thank you for registering an account with us. To complete your registration and unlock all features, please verify your email address by clicking the button below:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #007BFF; color: white; padding: 15px 25px; text-decoration: none; border-radius: 5px; font-size: 16px;">Verify Email</a>
			</div>
			<p style="color: #555;">If you did not create an account, no further action is required.</p>
			<p style="color: #555;">If the button above does not work, please copy and paste the following link into your web browser:</p>
			<p style="word-break: break-all; color: #007BFF;"><a href="%s">%s</a></p>
			<p style="color: #555;">This verification link will expire in 24 hours.</p>
			<br/>
			<p style="color: #555;">Best regards,<br/>The GOeShop Team</p>
		</div>
	`, verificationLink, verificationLink, verificationLink)

	msg.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(msg)
}
