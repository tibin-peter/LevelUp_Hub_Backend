package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOTPEmail(toEmail, otp string) error {
	from := mail.NewEmail("LevelUp Hub", os.Getenv("EMAIL_FROM"))
	to := mail.NewEmail("User", toEmail)

	subject := "Your OTP Code"

	plainText := fmt.Sprintf("Your OTP is:%s", otp)

	htmlText := fmt.Sprintf(`
<h2>LevelUp Hub Verification</h2>
<p>Hello,</p>
<p>Your verification code for LevelUp Hub is:</p>

<h1 style="letter-spacing:3px;">%s</h1>

<p>This code expires in 5 minutes.</p>
<p>If you didn’t request this, ignore this email.</p>

<br>
<p>— LevelUp Hub Team</p>
`, otp)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlText)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)

	if err != nil {
		fmt.Println("SenderGrid error:0", err)
		return err
	}

	fmt.Println("SendGrid status:", response.StatusCode)
	fmt.Println("SendGrid body:", response.Body)

	return nil
}
