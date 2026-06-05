package pkg

import (
	"crypto/tls"
	"net/smtp"
	"os"
)

func SendMail(receivers []string, subject, body string) error {
	host := "smtp.sumopod.com"
	port := "465"
	from := "support@viketin.id"
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	address := host + ":" + port

	message := []byte(
		"From: Tickitz <support@viketin.id>\r\n" +
			"To: " + receivers[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	auth := smtp.PlainAuth("", username, password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", address, tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range receivers {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(message)
	return err
}
