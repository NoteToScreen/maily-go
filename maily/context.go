package maily

import (
	"bytes"
	"math/rand"
	"mime/quotedprintable"
	"net/smtp"
	"strconv"
	"time"
)

// A Context contains information on how Maily should send emails. It is goroutine-safe. (that is, it may be accessed from multiple goroutines at the same time)
type Context struct {
	// The email address to send from
	FromAddress string

	// The from name, for example: "John Smith <jsmith@example.com>"
	FromDisplay string

	// The domain emails are being sent from (used to generate a Message-ID)
	SendDomain string

	// SMTP credentials, used to send the email
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// The path to the template folder. Can be relative or absolute.
	TemplatePath string
}

// An EmailResult struct contains information about an email that was sent successfully.
type EmailResult struct {
	// The value of the Message-ID header.
	MessageID string `json:"messageID"`
}

// SendMail sends an email to the given email address, using the given template and data.
func (c *Context) SendMail(toName string, toEmail string, templateName string, data TemplateData, textFuncs FuncMap, htmlFuncs FuncMap) (EmailResult, error) {
	if c.FromDisplay == "" || c.SendDomain == "" {
		return EmailResult{}, ErrInvalidConfig
	}

	subject, textBody, htmlBody, err := c.render(toEmail, templateName, data, textFuncs, htmlFuncs)
	if err != nil {
		return EmailResult{}, err
	}

	textEncodedBuf := bytes.NewBufferString("")
	textEncodedWriter := quotedprintable.NewWriter(textEncodedBuf)
	_, err = textEncodedWriter.Write([]byte(textBody))
	if err != nil {
		textEncodedWriter.Close()
		return EmailResult{}, err
	}
	textEncodedWriter.Close()

	htmlEncodedBuf := bytes.NewBufferString("")
	htmlEncodedWriter := quotedprintable.NewWriter(htmlEncodedBuf)
	_, err = htmlEncodedWriter.Write([]byte(htmlBody))
	if err != nil {
		htmlEncodedWriter.Close()
		return EmailResult{}, err
	}
	htmlEncodedWriter.Close()

	messageIDRandom := strconv.FormatInt(time.Now().Unix(), 10) + "." + strconv.Itoa(rand.Intn(999999))
	messageID := messageIDRandom + "@" + c.SendDomain

	toFull := toEmail
	if toName != "" {
		toFull = toName + " <" + toEmail + ">"
	}

	rawMessage := []byte(
		"From: " + c.FromDisplay + "\r\n" +
			"To: " + toFull + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Sender: " + c.FromDisplay + "\r\n" +
			"Message-ID: <" + messageID + ">\r\n" +
			"Date: " + time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700") + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: multipart/alternative; boundary=\"mimeboundary\"\r\n\r\n" +
			"--mimeboundary\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: quoted-printable\r\n\r\n" +
			textEncodedBuf.String() + "\r\n" +
			"--mimeboundary\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: quoted-printable\r\n\r\n" +
			htmlEncodedBuf.String() + "\r\n" +
			"--mimeboundary--")

	auth := smtp.PlainAuth("", c.SMTPUsername, c.SMTPPassword, c.SMTPHost)

	err = smtp.SendMail(c.SMTPHost+":"+strconv.Itoa(c.SMTPPort), auth, c.FromAddress, []string{toEmail}, rawMessage)
	if err != nil {
		return EmailResult{}, err
	}

	return EmailResult{
		MessageID: messageID,
	}, nil
}
