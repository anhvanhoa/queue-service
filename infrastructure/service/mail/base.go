package mail

import (
	"log"
	mailI "queue-service/domain/service/mail"
	"time"

	"github.com/wneessen/go-mail"
)

type mailProvider struct {
	mail     *mail.Msg
	provider *mail.Client
	config   *mailI.ConfigMail
}

func (m *mailProvider) SendMail(to []string, subject, body string, data map[string]any) error {
	m.mail.SetGenHeader("Content-Type", "text/html")
	m.mail.SetGenHeader("charset", "utf-8")
	m.mail.SetGenHeader("Date", time.Now().Format(time.RFC1123Z))
	m.mail.Subject(subject)
	m.mail.AddAlternativeString(mail.TypeTextHTML, body)
	m.mail.FromFormat(m.config.Name, m.config.Email)
	m.mail.To(to...)
	if err := m.provider.DialAndSend(m.mail); err != nil {
		return err
	}
	m.mail.Reset()
	return nil
}

func (m *mailProvider) SetProvider(cf *mailI.ConfigMail) mailI.MailProvider {
	opts := []mail.Option{
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithPort(cf.Port),
		mail.WithUsername(cf.UserName),
		mail.WithPassword(cf.Password),
	}
	if cf.TSL != nil {
		opts = append(opts, mail.WithTLSConfig(cf.TSL))
	}

	provider, err := mail.NewClient(cf.Host, opts...)
	if err != nil {
		log.Fatalf("Failed to create mail provider: %v", err)
	}
	m.provider = provider
	m.config = cf
	return m
}

func NewMailProvider() mailI.MailProvider {
	return &mailProvider{
		mail: mail.NewMsg(),
	}
}
