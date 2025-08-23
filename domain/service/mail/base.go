package mail

import "crypto/tls"

type ConfigMail struct {
	Host     string
	Port     int
	UserName string
	Password string
	Email    string
	Name     string
	TSL      *tls.Config
}

type MailProvider interface {
	SetProvider(cf *ConfigMail) MailProvider
	SendMail(to []string, subject, body string, data map[string]any) error
}
