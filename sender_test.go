package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const EMAIL_SENDER_NAME = "Example Name"
const EMAIL_SENDER_ADRESS = "examplemail@gmail.com"
const EMAIL_SENDER_PASSWORD = "password that google give you"

func TestNewGmailSender(t *testing.T) {

	sender := NewGmailSender(EMAIL_SENDER_NAME, EMAIL_SENDER_ADRESS, EMAIL_SENDER_PASSWORD)

	subject := "Test subject"

	content := `
	<h1>EXAMPLE</h1>
	<p>Hello world!</p>
	`

	to := []string{"exampledestination@gmail.com"}

	attachFiles := []string{"rute to your file to send"}

	err := sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)

}
