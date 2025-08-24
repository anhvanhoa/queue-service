package constants

type StatusSentMail string

const (
	STATUS_SENT_MAIL_PENDING StatusSentMail = "pending"
	STATUS_SENT_MAIL_SENT    StatusSentMail = "sent"
	STATUS_SENT_MAIL_FAILED  StatusSentMail = "failed"
)

func (s StatusSentMail) String() string {
	return string(s)
}
